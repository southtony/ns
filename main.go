package main

import (
	"fmt"
	"sync"
	"time"
)

var taskStore *store

var stopSignalCh chan interface{}
var stopFlowCh chan interface{}

const (
	CONSUMER_AMOUNT = 5
	PRODUCER_AMOUNT = 5
	BUFFER_SIZE     = 10
)

func main() {
	// Simple mock of an external data store
	taskStore = &store{
		mu: &sync.Mutex{},
	}

	stopSignalCh = make(chan interface{}, 1)
	stopFlowCh = make(chan interface{})

	tasksToExec := make(chan *task, BUFFER_SIZE)

	for i := 0; i < PRODUCER_AMOUNT; i++ {
		go producer(tasksToExec)
	}

	wg := &sync.WaitGroup{}

	for i := 0; i < CONSUMER_AMOUNT; i++ {
		go consumer(tasksToExec, wg)
		wg.Add(1)
	}

	// It starts a routine which should handle a stop signal from
	// `stopSignalCh` and then notify all producers and consumers to stop
	// via `stopFlowCh` channel

	go stopFlowManager()

	// Sleep for wait some work
	time.Sleep(time.Second * 10)

	// Stop all workers here
	stopSignalCh <- true

	// Wait for all consumers
	wg.Wait()

	fmt.Println("Success: ")
	for _, t := range taskStore.success {
		fmt.Println(t.result, t.executionTime)
	}

	fmt.Println("Failed: ")
	for _, t := range taskStore.failed {
		fmt.Println(t.result, t.executionTime)
	}

}

func stopFlowManager() {
	<-stopSignalCh
	close(stopFlowCh)
}

func producer(tasksToExec chan<- *task) {
	for {
		task := newTask()

		if time.Now().Nanosecond()%2 > 0 {
			task.failed = true
			task.errorMsg = "Some error occured"
		}

		// Handle stop signal as early as possible
		select {
		case <-stopFlowCh:
			return
		default:
		}

		select {
		case <-stopFlowCh:
			return
		case tasksToExec <- task:
		}
	}
}

func consumer(tasksToExec <-chan *task, wg *sync.WaitGroup) {
	for {
		select {
		case <-stopFlowCh:
			wg.Done()
			return
		default:
		}

		select {
		case <-stopFlowCh:
			wg.Done()
			return
		case task := <-tasksToExec:
			startTime := time.Now()

			task.result = "task has been successed"

			pastTime := time.Now().Add(-20 * time.Second)

			if !task.isCreationTimeAheadOfTime(pastTime) {
				task.result = "something went wrong"
				task.failed = true
			}

			time.Sleep(time.Millisecond * 150)

			endTime := time.Now()
			task.executionTime = endTime.Sub(startTime)

			taskStore.insertTask(task)
		}

	}
}

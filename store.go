package main

import (
	"sync"
)

type store struct {
	failed  []*task
	success []*task
	mu      *sync.Mutex
}

func (s *store) insertTask(task *task) {
	s.mu.Lock()

	defer s.mu.Unlock()

	if task.failed {
		s.failed = append(s.failed, task)
	} else {
		s.success = append(s.success, task)
	}
}

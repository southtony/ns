package main

import (
	"time"
)

type task struct {
	id            int
	createdAt     time.Time
	executionTime time.Duration
	result        string
	failed        bool
	errorMsg      string
}

func newTask() *task {
	now := time.Now()

	return &task{
		id:        int(now.Unix()),
		createdAt: now,
		failed:    false,
	}
}

func (t task) isCreationTimeAheadOfTime(compareTime time.Time) bool {
	return t.createdAt.After(compareTime)
}

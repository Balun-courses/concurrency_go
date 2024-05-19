package main

import (
	"sync"
	"time"
)

type Scheduler struct {
	mutex  sync.Mutex
	tasks  map[int]*time.Timer
	closed bool
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks: make(map[int]*time.Timer),
	}
}

// SetTimeout run some function after some timeout
func (s *Scheduler) SetTimeout(key int, delay time.Duration, action func()) {
	if delay < 0 || action == nil {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return
	}

	if task, found := s.tasks[key]; found {
		task.Stop()
	}

	s.tasks[key] = time.AfterFunc(delay, action)
}

// CancelTimeout cancel running of some function
func (s *Scheduler) CancelTimeout(key int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return
	}

	if task, found := s.tasks[key]; found {
		task.Stop()
		delete(s.tasks, key)
	}
}

// Shutdown cancel all functions
func (s *Scheduler) Shutdown() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.closed = true
	for key, task := range s.tasks {
		task.Stop()
		delete(s.tasks, key)
	}
}

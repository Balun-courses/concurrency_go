package main

import (
	"sync"
	"time"
)

type Scheduler struct {
	tasks map[int]*time.Timer
	mutex sync.Mutex
}

func (s *Scheduler) Delay(key int, function func(), delay time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if task, found := s.tasks[key]; found {
		task.Stop()
	}

	if function == nil {
		delete(s.tasks, key)
	}

	s.tasks[key] = time.AfterFunc(delay, func() {
		function()

		s.mutex.Lock()
		delete(s.tasks, key)
		s.mutex.Unlock()
	})
}

func (s *Scheduler) Shutdown() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, task := range s.tasks {
		task.Stop()
	}
}

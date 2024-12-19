package main

import (
	"errors"
	"sync"
	"time"
)

type Scheduler struct {
	mutex   sync.Mutex
	closed  bool
	actions map[int]*time.Timer
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		actions: make(map[int]*time.Timer),
	}
}

// SetTimeout run some function after some timeout
func (s *Scheduler) SetTimeout(key int, delay time.Duration, action func()) error {
	if delay < 0 {
		return errors.New("invalid delay")
	}

	if action == nil {
		return errors.New("invalid action")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.closed {
		return errors.New("scheduler is closed")
	}

	if timer, found := s.actions[key]; found {
		timer.Stop()
	}

	s.actions[key] = time.AfterFunc(delay, action)
	return nil
}

// CancelTimeout cancel running of some function
func (s *Scheduler) CancelTimeout(key int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if timer, found := s.actions[key]; found {
		timer.Stop()
		delete(s.actions, key)
	}
}

// Close cancel all functions
func (s *Scheduler) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.closed = true
	for key, timer := range s.actions {
		timer.Stop()
		delete(s.actions, key)
	}
}

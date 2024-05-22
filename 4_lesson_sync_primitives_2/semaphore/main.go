package main

import (
	"sync"
)

type Semaphore struct {
	count     int
	max       int
	condition *sync.Cond
}

func NewSemaphore(limit int) *Semaphore {
	mutex := &sync.Mutex{}
	return &Semaphore{
		max:       limit,
		condition: sync.NewCond(mutex),
	}
}

func (s *Semaphore) Acquire() {
	s.condition.L.Lock()
	defer s.condition.L.Unlock()

	for s.count >= s.max {
		s.condition.Wait()
	}

	s.count++
}

func (s *Semaphore) Release() {
	s.condition.L.Lock()
	defer s.condition.L.Unlock()

	s.count--
	s.condition.Signal()
}

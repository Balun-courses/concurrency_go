package main

import (
	"sync"
)

type Semaphore struct {
	count     int
	max       int
	mutex     *sync.Mutex
	condition *sync.Cond
}

func NewSemaphore(limit int) *Semaphore {
	mutex := &sync.Mutex{}
	return &Semaphore{
		max:       limit,
		mutex:     mutex,
		condition: sync.NewCond(mutex),
	}
}

func (s *Semaphore) Acquire() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for s.count >= s.max {
		s.condition.Wait()
	}

	s.count++
}

func (s *Semaphore) Release() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.count--
	s.condition.Signal()
}

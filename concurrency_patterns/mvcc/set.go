package main

import "sync"

type Set struct {
	mutex sync.RWMutex
	data  map[int64]struct{}
}

func NewSet() *Set {
	return &Set{
		data: make(map[int64]struct{}),
	}
}

func (s *Set) Insert(id int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[id] = struct{}{}
}

func (s *Set) Delete(id int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, id)
}

func (s *Set) Contains(id int64) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, found := s.data[id]
	return found
}

package main

import "sync"

type InMemoryStorage struct {
	mutex sync.RWMutex
	data  map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

func (s *InMemoryStorage) Set(key string, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[key] = value
}

func (s *InMemoryStorage) Get(key string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.data[key]
}

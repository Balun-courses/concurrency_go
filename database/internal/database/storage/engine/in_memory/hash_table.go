package in_memory

import (
	"sync"
)

var HashTableBuilder = func() hashTable {
	return NewHashTable()
}

type HashTable struct {
	mutex sync.RWMutex
	data  map[string]string
}

func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]string),
	}
}

func (s *HashTable) Set(key, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = value
}

func (s *HashTable) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, found := s.data[key]
	return value, found
}

func (s *HashTable) Del(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, key)
}

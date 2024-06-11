package main

import (
	"github.com/igrmk/treemap/v2"
	"sync"
)

type VersionedKey struct {
	Key  string
	TxID int32
}

func Less(lhs, rhs VersionedKey) bool {
	if lhs.Key < rhs.Key {
		return true
	} else if lhs.Key == rhs.Key {
		return lhs.TxID < rhs.TxID
	} else {
		return false
	}
}

type InMemoryStorage struct {
	mutex sync.RWMutex
	data  *treemap.TreeMap[VersionedKey, string]
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: treemap.NewWithKeyCompare[VersionedKey, string](Less),
	}
}

func (s *InMemoryStorage) Set(txID int32, modified map[string]string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for key, value := range modified {
		versionedKey := VersionedKey{
			Key:  key,
			TxID: txID,
		}

		s.data.Set(versionedKey, value)
	}
}

func (s *InMemoryStorage) ExistsBetween(readTxID, writeTxID int32, modified map[string]string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for key, _ := range modified {
		lowerBound := VersionedKey{
			Key:  key,
			TxID: readTxID,
		}

		upperBound := VersionedKey{
			Key:  key,
			TxID: writeTxID,
		}

		if begin, end := s.data.Range(lowerBound, upperBound); begin != end {
			return true
		}
	}

	return false
}

func (s *InMemoryStorage) Get(txID int32, key string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	lowerBound := VersionedKey{
		Key: key,
	}

	upperBound := VersionedKey{
		Key:  key,
		TxID: txID,
	}

	var value string
	begin, end := s.data.Range(lowerBound, upperBound)
	for ; begin != end; begin.Next() { // reverse
		value = begin.Value()
	}

	return value
}

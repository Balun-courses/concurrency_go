package main

import (
	"sync"
	"sync/atomic"
)

type storage interface {
	Set(string, string)
	Get(string) string
}

type txLock struct {
	kind  int32
	txId  int32
	mutex sync.RWMutex
}

type Scheduler struct {
	mutex      sync.Mutex
	locks      map[string]*txLock
	storage    storage
	identifier int32
}

func NewScheduler(storage storage) *Scheduler {
	return &Scheduler{
		locks:   make(map[string]*txLock),
		storage: storage,
	}
}

func (s *Scheduler) StartTransaction() Transaction {
	txId := atomic.AddInt32(&s.identifier, 1)
	return newTransaction(s, txId)
}

func (s *Scheduler) set(txId int32, key, value string) txOperation {
	var lock *txLock
	withLock(&s.mutex, func() {
		var found bool
		lock, found = s.locks[key]
		if !found {
			lock = &txLock{}
			s.locks[key] = lock
		}
	})

	operation := txOperation{
		key:  key,
		lock: lock,
	}

	if atomic.LoadInt32(&lock.txId) != txId {
		lock.mutex.Lock()
		atomic.StoreInt32(&lock.txId, txId)
		atomic.StoreInt32(&lock.kind, exclusive)
		previousValue := s.storage.Get(key)
		operation.value = &previousValue
	} else if atomic.LoadInt32(&lock.kind) == shared {
		lock.mutex.RUnlock()
		// unsafe with rollback
		lock.mutex.Lock()
		atomic.StoreInt32(&lock.kind, exclusive)
		previousValue := s.storage.Get(key)
		operation.value = &previousValue
	}

	s.storage.Set(key, value)
	return operation
}

func (s *Scheduler) get(txId int32, key string) (string, txOperation) {
	var lock *txLock
	withLock(&s.mutex, func() {
		var found bool
		lock, found = s.locks[key]
		if !found {
			lock = &txLock{}
			s.locks[key] = lock
		}
	})

	if atomic.LoadInt32(&lock.txId) != txId {
		lock.mutex.RLock()
		atomic.StoreInt32(&lock.txId, txId)
		atomic.StoreInt32(&lock.kind, shared)
	}

	operation := txOperation{
		key:  key,
		lock: lock,
	}

	return s.storage.Get(key), operation
}

func (s *Scheduler) commit(operations []txOperation) {
	const withoutRollback = false
	s.apply(operations, withoutRollback)
}

func (s *Scheduler) rollback(operations []txOperation) {
	const withRollback = true
	s.apply(operations, withRollback)
}

func (s *Scheduler) apply(operations []txOperation, rollback bool) {
	release := func(unlocked map[string]struct{}, operation txOperation, mutex sync.Locker) {
		atomic.StoreInt32(&operation.lock.txId, 0)
		atomic.StoreInt32(&operation.lock.kind, 0)
		unlocked[operation.key] = struct{}{}
		mutex.Unlock()
	}

	unlocked := make(map[string]struct{})
	for i := len(operations) - 1; i >= 0; i-- {
		operation := operations[i]
		_, alreadyUnlocked := unlocked[operation.key]
		if atomic.LoadInt32(&operation.lock.kind) == shared && !alreadyUnlocked {
			release(unlocked, operation, operation.lock.mutex.RLocker())
		} else if atomic.LoadInt32(&operation.lock.kind) == exclusive && operation.value != nil {
			if rollback {
				s.storage.Set(operation.key, *operation.value)
			}

			release(unlocked, operation, &operation.lock.mutex)
		}
	}
}

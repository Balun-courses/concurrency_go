package main

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type storage interface {
	Set(int32, map[string]string)
	Get(int32, string) string
	ExistsBetween(int32, int32, map[string]string) bool
}

type lock struct {
	locked atomic.Int32
}

type Scheduler struct {
	storage    storage
	identifier int32

	mutex sync.Mutex
	locks map[string]*lock
}

func NewScheduler(storage storage) *Scheduler {
	return &Scheduler{
		storage: storage,
		locks:   make(map[string]*lock),
	}
}

func (s *Scheduler) StartTransaction() Transaction {
	txId := atomic.AddInt32(&s.identifier, 1)
	return newTransaction(s, txId)
}

func (s *Scheduler) get(txId int32, key string) string {
	var l *lock
	withLock(&s.mutex, func() {
		l = s.locks[key]
	})

	if l != nil {
		ownerTxID := l.locked.Load()
		for ownerTxID != 0 && ownerTxID < s.identifier {
			runtime.Gosched()
			ownerTxID = l.locked.Load()
		}
	}

	return s.storage.Get(txId, key)
}

// Only "Snapshot Isolation" without "Serializable Snapshot Isolation"
func (s *Scheduler) commit(readTxId int32, modified map[string]string) bool {
	var locks []*lock
	withLock(&s.mutex, func() {
		for key, _ := range modified {
			l := s.locks[key]
			if l == nil {
				l = &lock{}
				s.locks[key] = l
			}

			locks = append(locks, l)
		}
	})

	if !s.acquireLocks(locks, readTxId) {
		return false
	}

	defer s.releaseLocks(locks)

	writeTxID := atomic.AddInt32(&s.identifier, 1)
	if s.storage.ExistsBetween(readTxId, writeTxID, modified) {
		return false
	}

	s.storage.Set(writeTxID, modified)
	return true
}

func (s *Scheduler) acquireLocks(locks []*lock, txID int32) bool {
	for idx, l := range locks {
		if !l.locked.CompareAndSwap(0, txID) {
			for idx >= 0 {
				locks[idx].locked.Store(0)
			}

			return false
		}
	}

	return true
}

func (s *Scheduler) releaseLocks(locks []*lock) {
	for _, l := range locks {
		l.locked.Store(0)
	}
}

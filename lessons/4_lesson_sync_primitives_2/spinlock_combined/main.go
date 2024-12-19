package main

import (
	"runtime"
	"sync/atomic"
)

const (
	unlocked = false
	locked   = true
)

const retriesNumber = 3

type SpinLock struct {
	state atomic.Bool
}

func NewSpinLock() *SpinLock {
	return &SpinLock{}
}

func (s *SpinLock) Lock() {
	retries := retriesNumber
	for !s.state.CompareAndSwap(unlocked, locked) {
		retries--
		if retries == 0 {
			runtime.Gosched()
			retries = retriesNumber
		}
	}
}

func (s *SpinLock) Unlock() {
	s.state.Store(unlocked)
}

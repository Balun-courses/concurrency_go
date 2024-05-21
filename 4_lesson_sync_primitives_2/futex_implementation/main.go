package main

import "sync/atomic"

type SpinLock struct {
	state atomic.Bool
}

func NewSpinLock() *SpinLock {
	return &SpinLock{}
}

func (s *SpinLock) Lock() {
	retries := 5
	for !s.state.CompareAndSwap(false, true) {
		retries--
		if retries == 0 {
			// park goroutine
			retries = 5
		}
	}
}

func (s *SpinLock) Unlock() {
	if !s.state.CompareAndSwap(true, false) {
		panic("incorrect usage")
	}
}

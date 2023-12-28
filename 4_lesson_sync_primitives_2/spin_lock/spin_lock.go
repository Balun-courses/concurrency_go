package main

import (
	"sync/atomic"
)

type SpinLock struct {
	state atomic.Bool
}

func NewSpinLock() SpinLock {
	return SpinLock{}
}

func (s *SpinLock) Lock() {
	for !s.state.CompareAndSwap(false, true) {
	}
}

func (s *SpinLock) TryLock() bool {
	return s.state.CompareAndSwap(false, true)
}

func (s *SpinLock) Unlock() {
	if !s.state.CompareAndSwap(true, false) {
		panic("incorrect ussage")
	}
}

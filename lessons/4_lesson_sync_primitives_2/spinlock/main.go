package main

import (
	"sync/atomic"
)

const (
	unlocked = false
	locked   = true
)

type SpinLock struct {
	state atomic.Bool
}

func NewSpinLock() *SpinLock {
	return &SpinLock{}
}

func (s *SpinLock) Lock() {
	for !s.state.CompareAndSwap(unlocked, locked) {
		// итерация за итерацией...
	}
}

func (s *SpinLock) Unlock() {
	s.state.Store(unlocked)
}

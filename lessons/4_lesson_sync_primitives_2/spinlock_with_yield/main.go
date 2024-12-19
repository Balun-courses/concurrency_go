package main

import (
	"runtime"
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
		runtime.Gosched() // но горутина не перейдет в состояние ожидания
	}
}

func (s *SpinLock) Unlock() {
	s.state.Store(unlocked)
}

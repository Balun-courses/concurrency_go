package main

import (
	"runtime"
	"sync/atomic"
)

const (
	unlocked = false
	locked   = true
)

type BrokenMutex struct {
	want   [2]atomic.Bool
	victim atomic.Int32
	owner  int
}

func (m *BrokenMutex) Lock(index int) {
	m.want[index].Store(locked)
	m.victim.Store(int32(index))

	for m.want[1-index].Load() && m.victim.Load() == int32(index) {
		runtime.Gosched()
	}

	m.owner = index
}

func (m *BrokenMutex) Unlock(index int) {
	m.want[m.owner].Store(unlocked)
}

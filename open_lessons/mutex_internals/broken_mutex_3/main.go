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
	want  [2]atomic.Bool
	owner int
}

// Здесь нет гарантии прогресса (liveness), так как могут несколько
// горутин могут помешать друг другу и попасть в livelock

func (m *BrokenMutex) Lock(index int) {
	m.want[index].Store(locked)
	for m.want[1-index].Load() {
		runtime.Gosched()
	}

	m.owner = index
}

func (m *BrokenMutex) Unlock(index int) {
	m.want[m.owner].Store(unlocked)
}

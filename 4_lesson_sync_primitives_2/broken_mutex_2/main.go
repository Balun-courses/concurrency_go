package main

import "sync/atomic"

const (
	unlocked = false
	locked   = true
)

type BrokenMutex struct {
	state atomic.Bool
}

// Здесь нет гарантии взаимного исключения (safety), так как несколько
// горутин могут попасть совместно в критическую секцию

func (m *BrokenMutex) Lock() {
	for m.state.Load() {
		// итерация за итерацией...
	}

	m.state.Store(locked)
}

func (m *BrokenMutex) Unlock() {
	m.state.Store(unlocked)
}

package main

import (
	"runtime"
	"sync"
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
	otherIndex := 1 - index
	m.want[index].Store(locked)
	for m.want[otherIndex].Load() {
		runtime.Gosched()
	}

	m.owner = index
}

func (m *BrokenMutex) Unlock() {
	m.want[m.owner].Store(unlocked)
}

const goroutinesNumber = 2

func main() {
	var mutex BrokenMutex
	wg := sync.WaitGroup{}
	wg.Add(goroutinesNumber)

	go func() {
		defer wg.Done()

		const goroutineIdx = 0
		mutex.Lock(goroutineIdx)
		mutex.Unlock()
	}()

	go func() {
		defer wg.Done()

		const goroutineIdx = 1
		mutex.Lock(goroutineIdx)
		mutex.Unlock()
	}()

	wg.Wait()
}

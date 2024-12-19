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

type Mutex struct {
	want   [2]atomic.Bool
	victim atomic.Int32
	owner  int
}

func (m *Mutex) Lock(index int) {
	m.want[index].Store(locked)
	m.victim.Store(int32(index))

	otherIndex := 1 - index
	for m.want[otherIndex].Load() && m.victim.Load() == int32(index) {
		runtime.Gosched()
	}

	m.owner = index
}

func (m *Mutex) Unlock() {
	m.want[m.owner].Store(unlocked)
}

const goroutinesNumber = 2

func main() {
	var mutex Mutex
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

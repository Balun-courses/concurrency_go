package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

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
		// iteration by iteration...
	}

	m.state.Store(locked)
}

func (m *BrokenMutex) Unlock() {
	m.state.Store(unlocked)
}

const goroutinesNumber = 1000

func main() {
	var mutex BrokenMutex
	wg := sync.WaitGroup{}
	wg.Add(goroutinesNumber)

	value := 0
	for i := 0; i < goroutinesNumber; i++ {
		go func() {
			defer wg.Done()

			mutex.Lock()
			value++
			mutex.Unlock()
		}()
	}

	wg.Wait()

	fmt.Println(value)
}

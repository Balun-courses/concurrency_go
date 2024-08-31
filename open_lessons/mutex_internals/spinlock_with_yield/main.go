package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

const (
	unlocked = false
	locked   = true
)

type Mutex struct {
	state atomic.Bool
}

func (m *Mutex) Lock() {
	for !m.state.CompareAndSwap(unlocked, locked) {
		runtime.Gosched() // но горутина не перейдет в состояние ожидания
	}
}

func (m *Mutex) Unlock() {
	m.state.Store(unlocked)
}

const goroutinesNumber = 1000

func main() {
	var mutex Mutex
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

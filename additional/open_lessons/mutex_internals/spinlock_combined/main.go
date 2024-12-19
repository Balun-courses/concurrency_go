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

const retriesNumber = 3

type Mutex struct {
	state atomic.Bool
}

func (m *Mutex) Lock() {
	retries := retriesNumber
	for !m.state.CompareAndSwap(unlocked, locked) {
		retries--
		if retries == 0 {
			runtime.Gosched()
			retries = retriesNumber
		}
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

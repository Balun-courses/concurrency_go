package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

type Mutex struct {
	ownerTicket    atomic.Int64
	nextFreeTicket atomic.Int64
}

func (m *Mutex) Lock() {
	ticket := m.nextFreeTicket.Add(1)
	for m.ownerTicket.Load() != ticket-1 {
		runtime.Gosched()
	}
}

func (m *Mutex) Unlock() {
	m.ownerTicket.Add(1)
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

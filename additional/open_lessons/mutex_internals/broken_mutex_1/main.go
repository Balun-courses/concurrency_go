package main

import (
	"fmt"
	"sync"
)

const (
	unlocked = false
	locked   = true
)

type BrokenMutex struct {
	state bool
}

// Здесь есть data race и нет гарантии взаимного исключения (safety),
// так как несколько горутин могут попасть совместно в критическую секцию

func (m *BrokenMutex) Lock() {
	for m.state {
		// iteration by iteration...
	}

	m.state = locked
}

func (m *BrokenMutex) Unlock() {
	m.state = unlocked
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

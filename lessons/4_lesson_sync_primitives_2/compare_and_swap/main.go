package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Need to show solution

var data map[string]string
var initialized atomic.Bool

func initialize() {
	if !initialized.Load() {
		initialized.Store(true)
		data = make(map[string]string)
		fmt.Println("initialized")
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			initialize()
		}()
	}

	wg.Wait()
}

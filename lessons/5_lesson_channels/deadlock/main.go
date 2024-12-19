package main

import (
	"sync"
)

var actions int
var mutex sync.Mutex
var buffer chan struct{}

func consumer() {
	for i := 0; i < 1000; i++ {
		mutex.Lock()
		actions++
		<-buffer
		mutex.Unlock()
	}
}

func producer() {
	for i := 0; i < 1000; i++ {
		buffer <- struct{}{}
		mutex.Lock()
		actions++
		mutex.Unlock()
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	buffer = make(chan struct{}, 1)

	go func() {
		defer wg.Done()
		consumer()
	}()

	go func() {
		defer wg.Done()
		producer()
	}()

	wg.Wait()
}

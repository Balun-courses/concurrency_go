package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Need to show solution

var mutex1 sync.Mutex
var mutex2 sync.Mutex

func goroutine1() {
	mutex1.Lock()

	runtime.Gosched()
	for !mutex2.TryLock() {
		// active waiting
	}

	mutex2.Unlock()
	mutex1.Unlock()

	fmt.Println("goroutine1 finished")
}

func goroutine2() {
	mutex2.Lock()

	runtime.Gosched()
	for !mutex1.TryLock() {
		// active waiting
	}

	mutex1.Unlock()
	mutex2.Unlock()

	fmt.Println("goroutine2 finished")
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		goroutine1()
	}()

	go func() {
		defer wg.Done()
		goroutine2()
	}()

	wg.Wait()
}

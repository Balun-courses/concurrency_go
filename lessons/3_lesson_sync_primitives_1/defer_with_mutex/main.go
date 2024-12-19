package main

import "sync"

// Need to show solution

var mutex sync.Mutex

func operation() {}

func doSomething() {
	mutex.Lock()
	operation()
	mutex.Unlock()

	// some long operation

	mutex.Lock()
	operation()
	mutex.Unlock()
}

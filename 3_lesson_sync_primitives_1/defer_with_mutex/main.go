package main

import "sync"

// Need to show solution

var mutex sync.Mutex

func operation() error {
	return nil // or error
}

func doSomething() {
	mutex.Lock()

	err := operation()
	if err != nil {
		mutex.Unlock()
		return
	}

	err = operation()
	if err != nil {
		mutex.Unlock()
		return
	}

	mutex.Unlock()
}

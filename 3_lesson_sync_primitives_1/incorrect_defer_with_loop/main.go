package main

import (
	"sync"
)

// Need to show solution

var mutex sync.Mutex
var values []int

func doSomething() {
	for number := 0; number < 10; number++ {
		mutex.Lock()
		defer mutex.Unlock()

		values = append(values, number)
	}
}

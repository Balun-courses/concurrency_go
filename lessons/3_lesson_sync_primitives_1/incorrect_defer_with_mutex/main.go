package main

import (
	"fmt"
	"sync"
)

// Need to show solution

var mutex sync.Mutex
var cache map[string]string

func doSomething() {
	var value string

	{
		mutex.Lock()
		defer mutex.Unlock()
		value = cache["key"]
	}

	fmt.Println(value)

	{
		mutex.Lock()
		defer mutex.Unlock()
		cache["key"] = "value"
	}
}

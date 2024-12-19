package main

import (
	"fmt"
	"sync"
)

var mutex sync.Mutex
var value string

func set(v string) {
	mutex.Lock()
	value = v
	mutex.Unlock()
}

func print() {
	mutex.Lock()
	fmt.Println(value)
	mutex.Unlock()
}

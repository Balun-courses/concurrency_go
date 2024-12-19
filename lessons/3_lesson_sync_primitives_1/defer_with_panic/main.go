package main

import (
	"log"
	"sync"
)

var mutex sync.Mutex

func functionWithPanic() {
	panic("error")
}

func handle1() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("recovered")
		}
	}()

	mutex.Lock()
	defer mutex.Unlock()

	functionWithPanic()
}

func handle2() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("recovered")
		}
	}()

	mutex.Lock()
	functionWithPanic()
	mutex.Unlock()
}

func main() {
	handle1()
	handle2()
}

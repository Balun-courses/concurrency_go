package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	var mutex sync.Mutex
	mutex.Lock()

	go func() {
		time.Sleep(time.Second)
		log.Println("Hi")
		mutex.Unlock()
	}()

	mutex.Lock()
	log.Println("Bye")
}

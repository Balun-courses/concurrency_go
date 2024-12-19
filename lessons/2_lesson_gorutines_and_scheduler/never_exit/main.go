package main

import (
	"log"
	"time"
)

func task() {
	for {
		time.Sleep(time.Millisecond * 200)
		panic("unexpected situation")
	}
}

func NeverExit(name string, action func()) {
	defer func() {
		if v := recover(); v != nil {
			log.Println(name, "is crashed - restarting...")
			go NeverExit(name, action)
		}
	}()

	if action != nil {
		action()
	}
}

func main() {
	go NeverExit("first_goroutine", task)
	go NeverExit("second_goroutine", task)

	time.Sleep(time.Second)
}

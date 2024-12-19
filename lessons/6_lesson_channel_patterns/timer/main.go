package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- struct{}) {
	time.Sleep(5 * time.Second)
	ch <- struct{}{}
}

func main() {
	ch := make(chan struct{})
	go producer(ch)

	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	for {
		select {
		case <-ch:
			return
		case <-timer.C:
			fmt.Println("tick")
		}
	}
}

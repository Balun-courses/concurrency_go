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

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ch:
			return
		case <-ticker.C:
			fmt.Println("tick")
		}
	}
}

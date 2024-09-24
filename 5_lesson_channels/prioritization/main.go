package main

import (
	"fmt"
	"time"
)

// Need to show solution

func producer(ch chan<- int) {
	for {
		ch <- 1
		time.Sleep(time.Second)
	}
}

func main() {
	ch1 := make(chan int) // more prioritized
	ch2 := make(chan int)

	go producer(ch1)
	go producer(ch2)

	for {
		select {
		case value := <-ch1:
			fmt.Println(value)
		case value := <-ch2:
			fmt.Println(value)
		}
	}
}

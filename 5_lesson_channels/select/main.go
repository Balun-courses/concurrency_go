package main

import (
	"fmt"
	"time"
)

func async1() chan string {
	ch := make(chan string)
	go func() {
		time.Sleep(1 * time.Second)
		ch <- "async1 result"
	}()
	return ch
}

func async2() chan string {
	ch := make(chan string)
	go func() {
		time.Sleep(1 * time.Second)
		ch <- "async2 result"
	}()
	return ch
}

func main() {
	ch1 := async1()
	ch2 := async2()

	select {
	case result := <-ch1:
		fmt.Println(result)
	case result := <-ch2:
		fmt.Println(result)
	}
}

package main

import (
	"fmt"
	"time"
)

// Need to show solution

func FetchData1() chan int {
	ch := make(chan int)
	go func() {
		time.Sleep(time.Second * 2)
		ch <- 10
	}()

	return ch
}

func FetchData2() chan int {
	ch := make(chan int)
	go func() {
		time.Sleep(time.Second * 2)
		ch <- 20
	}()

	return ch
}

func Process(value1, value2 int) {
	// Processing...
}

func main() {
	start := time.Now()
	Process(<-FetchData1(), <-FetchData2())
	fmt.Println(time.Now().Sub(start))
}

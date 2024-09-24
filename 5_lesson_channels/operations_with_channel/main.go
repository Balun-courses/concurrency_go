package main

import (
	"fmt"
)

func writeToNilChannel() {
	var ch chan int
	ch <- 1
}

func writeToClosedChannel() {
	ch := make(chan int, 2)
	close(ch)
	ch <- 20
}

func readFromChannel() {
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20

	val, ok := <-ch
	fmt.Println(val, ok)

	close(ch)
	val, ok = <-ch
	fmt.Println(val, ok)

	val, ok = <-ch
	fmt.Println(val, ok)
}

func readAnyChannels() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		ch1 <- 100
	}()

	go func() {
		ch2 <- 200
	}()

	select {
	case val1 := <-ch1:
		fmt.Println(val1)
	case val2 := <-ch2:
		fmt.Println(val2)
	}
}

func readFromNilChannel() {
	var ch chan int
	<-ch
}

func rangeNilChannel() {
	var ch chan int
	for range ch {

	}
}

func closeNilChannel() {
	var ch chan int
	close(ch)
}

func closeChannelAnyTimes() {
	ch := make(chan int)
	close(ch)
	close(ch)
}

func main() {
}

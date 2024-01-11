package main

import (
	"fmt"
	"time"
)

func writeToNilChannel() {
	var ch chan int
	ch <- 1
}

func readFromNilChannel() {
	var ch chan int
	<-ch
}

func closeNilChannel() {
	var ch chan int
	close(ch)
}

func rangeNilChannel() {
	var ch chan int
	for _ = range ch {

	}
}

func openNilChannel() {
	var ch chan int

	go func() {
		ch = make(chan int)
		ch <- 100
		close(ch)
	}()

	time.Sleep(100 * time.Millisecond)
	for value := range ch {
		fmt.Println(value)
	}
}

func closeChannelAnyTimes() {
	ch := make(chan int)
	close(ch)
	close(ch)
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

	// ch doesn't have data
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

func writeToClosedChannel() {
	ch := make(chan int, 2)
	ch <- 10

	close(ch)
	ch <- 20
}

func writeToClosedBufferedChannel() {
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20

	go func() {
		ch <- 30
	}()

	time.Sleep(100 * time.Millisecond)
	close(ch)

	for value := range ch {
		fmt.Println(value)
	}
}

func getEventAfterClose() {
	ch := make(chan int, 2)
	go func() {
		<-ch
		fmt.Println("event 1")
	}()

	time.Sleep(100 * time.Millisecond)
	close(ch)

	<-ch
	fmt.Println("event 2")
}

func main() {
}

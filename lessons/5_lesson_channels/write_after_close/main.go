package main

import (
	"time"
)

func main() {
	ch := make(chan int)
	go func() {
		ch <- 1
	}()

	time.Sleep(500 * time.Millisecond)

	close(ch)
	<-ch

	time.Sleep(100 * time.Millisecond)
}

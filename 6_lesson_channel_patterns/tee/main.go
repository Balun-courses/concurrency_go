package main

import (
	"fmt"
	"sync"
)

func Tee(in chan string) (chan string, chan string) {
	out1 := make(chan string)
	out2 := make(chan string)

	go func() {
		defer close(out1)
		defer close(out2)

		for value := range in {
			out1 <- value
			out2 <- value
		}
	}()

	return out1, out2
}

func main() {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for i := 0; i < 5; i++ {
			ch <- "test"
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	ch1, ch2 := Tee(ch)

	go func() {
		defer wg.Done()
		for value := range ch1 {
			fmt.Println("ch1: " + value)
		}
	}()

	go func() {
		defer wg.Done()
		for value := range ch2 {
			fmt.Println("ch2: " + value)
		}
	}()

	wg.Wait()
}

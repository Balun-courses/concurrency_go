package main

import (
	"fmt"
	"sync"
	"time"
)

func parse(in chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for data := range in {
			time.Sleep(50 * time.Millisecond)
			out <- fmt.Sprintf("parsed - %s", data)
		}
	}()

	return out
}

func send(in <-chan string) <-chan string {
	out := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for data := range in {
			time.Sleep(100 * time.Millisecond)
			out <- fmt.Sprintf("sent - %s", data)
		}
	}()

	go func() {
		defer wg.Done()
		for data := range in {
			time.Sleep(100 * time.Millisecond)
			out <- fmt.Sprintf("sent - %s", data)
		}
	}()

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for i := 0; i < 5; i++ {
			ch <- "value"
		}
	}()

	out := send(parse(ch))
	for value := range out {
		fmt.Println(value)
	}
}

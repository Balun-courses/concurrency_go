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
	splitChs := split(in, 2)
	out := make(chan string)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for data := range splitChs[0] {
			time.Sleep(100 * time.Millisecond)
			out <- fmt.Sprintf("sent - %s", data)
		}
	}()

	go func() {
		defer wg.Done()
		for data := range splitChs[1] {
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

func split(inputCh <-chan string, n int) []chan string {
	outputCh := make([]chan string, n)
	for i := 0; i < n; i++ {
		outputCh[i] = make(chan string)
	}

	go func() {
		idx := 0
		for value := range inputCh {
			outputCh[idx] <- value
			idx = (idx + 1) % n
		}

		for _, ch := range outputCh {
			close(ch)
		}
	}()

	return outputCh
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

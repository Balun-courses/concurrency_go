package main

import (
	"fmt"
	"time"
)

func OrDone(done chan struct{}, in chan string) chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case value, ok := <-in:
				if !ok {
					return
				}

				select {
				case out <- value:
				case <-done:
					return
				}
			}
		}
	}()

	return out
}

func main() {
	in := make(chan string)
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("test")
		}
	}()

	done := make(chan struct{})
	go func() {
		time.Sleep(time.Second)
		close(done)
	}()

	for value := range OrDone(done, in) {
		fmt.Println(value)
	}
}

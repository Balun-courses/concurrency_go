package main

import (
	"fmt"
	"time"
)

func Filter(input <-chan int) <-chan int {
	output := make(chan int)

	go func() {
		for number := range input {
			if number%2 != 0 {
				output <- number
			}
		}

		close(output)
	}()

	return output
}

func main() {
	in := make(chan int)

	go func() {
		defer close(in)
		for i := 0; i < 10; i++ {
			in <- i
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for value := range Filter(in) {
		fmt.Println(value)
	}
}

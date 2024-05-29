package main

import (
	"fmt"
)

func Bridge(in chan chan string) chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for ch := range in {
			for value := range ch {
				out <- value
			}
		}
	}()

	return out
}

func main() {
	in := make(chan chan string)
	go func() {
		innerCh1 := make(chan string, 3)
		for i := 0; i < 3; i++ {
			innerCh1 <- "inner-ch-1"
		}

		close(innerCh1)

		innerCh2 := make(chan string, 3)
		for i := 0; i < 3; i++ {
			innerCh2 <- "inner-ch-2"
		}

		close(innerCh2)

		in <- innerCh1
		in <- innerCh2
		close(in)
	}()

	for value := range Bridge(in) {
		fmt.Println(value)
	}
}

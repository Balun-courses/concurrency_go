package main

import "fmt"

func Bridge(in chan chan string) chan string {
	out := make(chan string)
	go func() {
		defer close(out)

		for {
			select {
			case ch, ok := <-in:
				if ok == false {
					return
				}

				for value := range ch {
					out <- value
				}
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

	out := Bridge(in)
	for value := range out {
		fmt.Println(value)
	}
}

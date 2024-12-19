package main

import "fmt"

func main() {
	ch1 := make(chan struct{}, 1)
	ch2 := make(chan struct{}, 1)

	close(ch1)
	close(ch2)

	ch1Value := 0.0
	ch2Value := 0.0

	for i := 0; i < 100000; i++ {
		select {
		case <-ch1:
			ch1Value++
		case <-ch1:
			ch1Value++
		case <-ch2:
			ch2Value++
		}
	}

	fmt.Println(ch1Value / ch2Value)
}

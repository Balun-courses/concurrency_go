package main

import "fmt"

func Transform(in <-chan int) <-chan int {
	result := make(chan int)
	go func() {
		defer close(result)
		for number := range in {
			result <- number * number
		}
	}()

	return result
}

func main() {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < 5; i++ {
			ch <- i
		}
	}()

	for number := range Transform(ch) {
		fmt.Println(number)
	}
}

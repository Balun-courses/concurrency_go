package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)

	value := 0
	go func() {
		ch <- 1
	}()
	go func() {
		ch <- 1
	}()

	value += <-ch
	value += <-ch

	fmt.Println(value)
}

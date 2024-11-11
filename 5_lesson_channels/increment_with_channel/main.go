package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)

	go func() {
		ch <- 1
	}()
	go func() {
		ch <- 1
	}()

	value := 0
	value += <-ch
	value += <-ch

	fmt.Println(value)
}

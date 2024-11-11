package main

import "fmt"

func main() {
	source := make(chan int)
	clone := source

	go func() {
		source <- 1
	}()

	fmt.Println(<-clone)
}

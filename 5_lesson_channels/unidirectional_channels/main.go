package main

import "fmt"

func in(in chan<- int) {
	in <- 100
	close(in)
}

func out(out <-chan int) {
	fmt.Println(<-out)
}

func main() {
	var ch = make(chan int, 1)
	in(ch)
	out(ch)
}

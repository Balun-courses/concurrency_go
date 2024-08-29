package main

import "fmt"

func main() {
	ch := make(chan int, 1)
	for done := false; !done; {
		select {
		default:
			fmt.Println(3)
			done = true
		case <-ch:
			fmt.Println(2)
			ch = nil
		case ch <- 1:
			fmt.Println(1)
		}
	}
}

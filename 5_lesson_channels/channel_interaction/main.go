package main

import (
	"fmt"
	"time"
)

func async(ch chan string) {
	time.Sleep(2 * time.Second)
	ch <- "async result"
}

func main() {
	ch := make(chan string)
	go async(ch)
	// ...
	result := <-ch
	fmt.Println(result)
}

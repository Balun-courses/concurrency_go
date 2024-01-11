package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(strings <-chan string) <-chan struct{} {
		completed := make(chan struct{})
		go func() {
			defer func() {
				fmt.Println("doWork exited")
				close(completed)
			}()

			for str := range strings {
				fmt.Println(str)
			}
		}()

		return completed
	}

	strings := make(chan string)
	doWork(strings)
	strings <- "Test"

	time.Sleep(time.Second)
	fmt.Println("Done")
}

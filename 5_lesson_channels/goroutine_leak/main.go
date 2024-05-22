package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(strings <-chan string) {
		go func() {
			for str := range strings {
				fmt.Println(str)
			}

			fmt.Println("doWork exited")
		}()
	}

	strings := make(chan string)
	doWork(strings)
	strings <- "Test"

	time.Sleep(time.Second)
	fmt.Println("Done")
}

package main

import (
	"fmt"
	"time"
)

// Need to show solution

func main() {
	for {
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("timeout")
			return
		case <-time.Tick(time.Second):
			fmt.Println("tick")
		}
	}
}

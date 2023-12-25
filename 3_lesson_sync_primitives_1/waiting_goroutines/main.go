package main

import "fmt"

// Need to show solution

func main() {
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Println("test")
		}()
	}
}

package main

import "log"

// Need to show solution

func main() {
	for i := 0; i < 5; i++ {
		go func() {
			log.Println("test")
		}()
	}
}

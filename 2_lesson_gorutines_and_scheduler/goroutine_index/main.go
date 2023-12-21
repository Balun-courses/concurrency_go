package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 5; i++ {
		go func() {
			fmt.Print(i)
		}()
	}

	time.Sleep(2 * time.Second)
}

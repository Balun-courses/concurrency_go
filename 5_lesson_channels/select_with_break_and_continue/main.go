package main

import "fmt"

func main() {
	data := make(chan int)
	go func() {
		for i := 1; i <= 4; i++ {
			data <- i
		}
		close(data)
	}()

	for {
		value := 0
		found := true

		select {
		case value, found = <-data:
			if value == 2 {
				continue
			} else if value == 3 {
				break
			}

			if !found {
				return
			}
		}

		fmt.Println(value)
	}
}

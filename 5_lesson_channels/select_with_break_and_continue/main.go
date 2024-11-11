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
		opened := true

		select {
		case value, opened = <-data:
			if value == 2 {
				continue
			} else if value == 3 {
				break
			}

			if !opened {
				return
			}
		}

		fmt.Println(value)
	}
}

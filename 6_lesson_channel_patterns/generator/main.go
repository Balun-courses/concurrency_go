package main

import "fmt"

func GenerateWithClosure(number int) func() int {
	return func() int {
		r := number
		number++
		return r
	}
}

func GenerateWithChannel(start, end int) chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for number := start; number <= end; number++ {
			ch <- number
		}
	}()

	return ch
}

func main() {
	generator := GenerateWithClosure(100)
	for i := 0; i <= 200; i++ {
		fmt.Println(generator())
	}

	for number := range GenerateWithChannel(100, 200) {
		fmt.Println(number)
	}
}

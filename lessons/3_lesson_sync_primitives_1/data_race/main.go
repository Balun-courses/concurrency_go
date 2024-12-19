package main

import (
	"fmt"
	"sync"
)

func main() {
	text := ""

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		text = "hello world"
	}()

	go func() {
		defer wg.Done()
		fmt.Println(text)
	}()

	wg.Wait()
}

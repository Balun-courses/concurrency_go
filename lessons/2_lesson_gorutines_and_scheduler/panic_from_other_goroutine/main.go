package main

import (
	"fmt"
	"time"
)

func process() {
	defer func() {
		v := recover()
		fmt.Println("recovered:", v)
	}()

	go func() {
		panic("error")
	}()

	time.Sleep(time.Second)
}

func main() {
	process()
}

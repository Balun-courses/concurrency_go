package main

import (
	"fmt"
	"runtime"
)

func infiniteLoop(str string) {
	fmt.Println("infinite loop")
	for {
		fmt.Println("test")
	}
}

func loop(str string) {
	for i := 0; i < 5; i++ {
		runtime.Gosched()
		fmt.Println(str)
	}
}

func main() {
	runtime.GOMAXPROCS(1)
	go infiniteLoop("infinite_loop")
	loop("loop")
}

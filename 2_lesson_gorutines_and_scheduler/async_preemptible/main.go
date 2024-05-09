package main

import (
	"fmt"
	"runtime"
)

func infiniteLoop(str string) {
	for {
		fmt.Println(str)
	}
}

func loop(str string) {
	for i := 0; i < 1; i++ {
		runtime.Gosched()
		fmt.Println(str)
	}
}

func main() {
	runtime.GOMAXPROCS(1)
	go infiniteLoop("infinite_loop")
	loop("loop")
}

package main

import (
	"log"
	"runtime"
)

func infiniteLoop(str string) {
	for {
		log.Println(str)
	}
}

func loop(str string) {
	for i := 0; i < 5; i++ {
		runtime.Gosched()
		log.Println(str)
	}
}

func main() {
	runtime.GOMAXPROCS(1)
	go infiniteLoop("infinite_loop")
	loop("loop")
}

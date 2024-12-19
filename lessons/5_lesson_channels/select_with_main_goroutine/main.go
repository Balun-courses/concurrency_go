package main

import "runtime"

func doSomething() {
	for {
		runtime.Gosched()
	}
}

func main() {
	go doSomething()
	go doSomething()
	select {}
}

package main

func f(out chan<- int) {
	out <- 100
}

func main() {
	var ch = make(chan int)
	f(ch)
}

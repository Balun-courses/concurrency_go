package main

func way1() {
	make(chan struct{}) <- struct{}{}
	// or
	make(chan<- struct{}) <- struct{}{}
}

func way2() {
	<-make(chan struct{})
	// or
	<-make(<-chan struct{})
	// or
	for range make(<-chan struct{}) {
	}
}

func way3() {
	chan struct{}(nil) <- struct{}{}
	// or
	<-chan struct{}(nil)
	// or
	for range chan struct{}(nil) {
	}
}

func way4() {
	select {}
}

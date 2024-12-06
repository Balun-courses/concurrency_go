package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	time.Sleep(time.Second * 2)
	fmt.Println("after sleep")

	/*
		Note, the if code block is used to discard/drain the potential
		ticker notification which is sent in the small period when executing
		the second branch code block (since Go 1.23, this has become needless).
		Note: since Go 1.23, the Ticker.Reset method will automatically
		discard/drain the potential stale ticker notification.

		https://go.dev/play/p/JAC8Ln1kwMz?v=goprev
	*/

	ticker.Reset(time.Second * 3)

	<-ticker.C
	fmt.Println("first")
	<-ticker.C
	fmt.Println("second")
}

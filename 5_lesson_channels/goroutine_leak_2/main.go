package main

// Need to show solution and describe close

// First-response-wins strategy
func request() int {
	ch := make(chan int)
	for i := 0; i < 5; i++ {
		go func() {
			ch <- i // 4 goroutines will be blocked
		}()
	}

	return <-ch
}

package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func receiveWeather(ctx context.Context, result chan struct{}, idx int) {
	randomTime := time.Duration(rand.Intn(5000)) * time.Millisecond

	timer := time.NewTimer(randomTime)
	defer timer.Stop()

	select {
	case <-timer.C:
		fmt.Printf("finished: %d\n", idx)
		result <- struct{}{}
	case <-ctx.Done():
		fmt.Printf("canceled: %d\n", idx)
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(10)

	ctx, cancel := context.WithCancel(context.Background())

	result := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func(idx int) {
			defer wg.Done()
			receiveWeather(ctx, result, idx)
		}(i)
	}

	<-result
	cancel()

	wg.Wait()
}

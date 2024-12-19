package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	makeRequest(ctx)
}

func makeRequest(ctx context.Context) {
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	newCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	select {
	case <-newCtx.Done():
		fmt.Println("canceled")
	case <-timer.C:
		fmt.Println("timer")
	}
}

package main

import (
	"context"
	"fmt"
	"time"
)

func makeRequest(ctx context.Context) {
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case <-timer.C:
		fmt.Println("finished")
	case <-ctx.Done():
		fmt.Println("canceled")
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	makeRequest(ctx)
}

package main

import (
	"context"
	"fmt"
	"time"
)

// context.AfterFunc

func WithAfterFunc(ctx context.Context, action func()) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	if action != nil {
		go func() {
			<-ctx.Done()
			action()
		}()
	}

	return ctx, cancel
}

func main() {
	afterDone := func() {
		fmt.Println("after")
	}

	_, cancel := WithAfterFunc(context.Background(), afterDone)
	cancel()

	time.Sleep(100 * time.Millisecond)
}

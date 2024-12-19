package main

import (
	"context"
	"fmt"
	"time"
)

// context.AfterFunc

func WithCtxAfterFunc(ctx context.Context, action func()) {
	if action != nil {
		go func() {
			<-ctx.Done()
			action()
		}()
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	WithCtxAfterFunc(ctx, func() {
		fmt.Println("after")
	})

	cancel()

	time.Sleep(100 * time.Millisecond)
}

package main

import (
	"context"
)

func WithAfterFunc(ctx context.Context, action func()) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	if action != nil {
		go func() {
			<-ctx.Done()
		}()
	}

	return ctx, cancel
}

func main() {
	afterDone := func() {
		// do some work
	}

	_, cancel := WithAfterFunc(context.Background(), afterDone)
	defer cancel()
}

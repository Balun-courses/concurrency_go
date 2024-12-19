package main

import "context"

func incorrectCheck(ctx context.Context, stream <-chan string) {
	data := <-stream
	_ = data
}

func correctCheck(ctx context.Context, stream <-chan string) {
	select {
	case data := <-stream:
		_ = data
	case <-ctx.Done():
		return
	}
}

package main

import "context"

func WithContexCheck(ctx context.Context, action func()) {
	if action == nil || ctx.Err() != nil {
		return
	}

	action()
}

func main() {
	ctx := context.Background()
	WithContexCheck(ctx, func() {
		// do something
	})
}

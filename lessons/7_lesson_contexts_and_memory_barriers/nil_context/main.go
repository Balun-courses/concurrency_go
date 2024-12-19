package main

import (
	"context"
)

func process(ctx context.Context) {
	if ctx.Err() != nil {
		// handling...
	}
}

func main() {
	process(nil)
}

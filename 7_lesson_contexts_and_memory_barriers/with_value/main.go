package main

import (
	"context"
	"fmt"
)

func main() {
	traceCtx := context.WithValue(context.Background(), "trace_id", "12-21-33")
	makeRequest(traceCtx)
}

func makeRequest(ctx context.Context) {
	oldValue, ok := ctx.Value("trace_id").(string)
	if ok {
		fmt.Println(oldValue)
	}

	newCtx := context.WithValue(ctx, "trace_id", "22-22-22")
	newValue, ok := newCtx.Value("trace_id").(string)
	if ok {
		fmt.Println(newValue)
	}
}

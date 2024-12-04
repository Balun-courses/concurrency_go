package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	innerCtx := context.WithoutCancel(ctx)
	cancel()

	if innerCtx.Err() != nil {
		fmt.Println("canceled")
	}
}

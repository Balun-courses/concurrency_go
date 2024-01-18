package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, cancel = context.WithCancel(ctx)
	cancel()

	select {
	case <-ctx.Done():
		fmt.Println("canceled")
	}
}

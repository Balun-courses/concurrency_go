package main

import (
	"context"
	"errors"
	"fmt"
)

func main() {
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(errors.New("error"))

	fmt.Println(ctx.Err())
	fmt.Println(context.Cause(ctx))
}

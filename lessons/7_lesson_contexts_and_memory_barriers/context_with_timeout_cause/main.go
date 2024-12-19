package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeoutCause(context.Background(), time.Second, errors.New("timeout"))
	defer cancel() // show difference

	<-ctx.Done()

	fmt.Println(ctx.Err())
	fmt.Println(context.Cause(ctx))
}

package main

import (
	"context"
	"errors"
	"fmt"
)

func main() {
	ctx, cancel := context.WithCancelCause(context.Background()) // added in go1.20
	cancel(errors.New("error"))

	if ctx.Err() != nil {
		err := context.Cause(ctx)
		fmt.Println(err.Error())
	}
}

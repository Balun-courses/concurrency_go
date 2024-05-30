package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, groupCtx := errgroup.WithContext(ctx)
	for i := 0; i < 10; i++ {
		group.Go(func() error {
			timeout := time.Second * time.Duration(rand.Intn(10))

			timer := time.NewTimer(timeout)
			defer timer.Stop()

			select {
			case <-timer.C:
				fmt.Println("timeout")
				return errors.New("error")
			case <-groupCtx.Done():
				fmt.Println("canceled")
				return nil
			}
		})
	}

	if err := group.Wait(); err != nil {
		fmt.Println(err.Error())
	}
}

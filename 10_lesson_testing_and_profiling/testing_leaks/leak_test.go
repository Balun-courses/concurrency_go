package testing_leaks

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"testing"
	"time"
)

func work(result chan int) {
	time.Sleep(2 * time.Second)
	result <- 100
}

func Manage() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result := make(chan int)
	go work(result)

	select {
	case <-result:
		return nil
	case <-ctx.Done():
		return errors.New("error")
	}
}

func TestGoroutineLeak(t *testing.T) {
	defer goleak.VerifyNone(t)

	err := Manage()
	require.Error(t, err, "error")
}

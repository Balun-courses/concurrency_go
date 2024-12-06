package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type ErrGroup struct {
	err    error
	wg     sync.WaitGroup
	once   sync.Once
	doneCh chan struct{}
}

func NewErrGroup() *ErrGroup {
	return &ErrGroup{
		doneCh: make(chan struct{}),
	}
}

func (eg *ErrGroup) Go(task func() error) {
	select {
	case <-eg.doneCh:
		return
	default:
	}

	eg.wg.Add(1)
	go func() {
		defer eg.wg.Done()

		select {
		case <-eg.doneCh:
			return
		default:
			if err := task(); err != nil {
				eg.once.Do(func() {
					eg.err = err
					close(eg.doneCh)
				})
			}
		}
	}()
}

func (eg *ErrGroup) Wait() error {
	eg.wg.Wait()
	return eg.err
}

func main() {
	group := NewErrGroup()
	group.Go(func() error {
		fmt.Println("started")
		return errors.New("error")
	})

	time.Sleep(time.Second)
	for i := 0; i < 5; i++ {
		group.Go(func() error {
			fmt.Println("started after timeout")
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		fmt.Println(err.Error())
	}
}

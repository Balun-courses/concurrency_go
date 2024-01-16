package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type ErrGroup struct {
	err    unsafe.Pointer
	wg     sync.WaitGroup
	doneCh chan struct{}
}

func NewErrGroup() *ErrGroup {
	return &ErrGroup{
		doneCh: make(chan struct{}),
	}
}

func (eg *ErrGroup) Go(task func() error) {
	select {
	case _, ok := <-eg.doneCh:
		if !ok {
			return
		}
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
				newPtr := unsafe.Pointer(&err)
				if atomic.CompareAndSwapPointer(&eg.err, nil, newPtr) {
					close(eg.doneCh)
				}
			}
		}
	}()
}

func (eg *ErrGroup) Wait() error {
	eg.wg.Wait()
	if err := atomic.LoadPointer(&eg.err); err != nil {
		return *(*error)(err)
	} else {
		return nil
	}
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

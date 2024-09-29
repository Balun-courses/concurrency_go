package main

import (
	"errors"
	"fmt"
	"sync"
)

type TaskPool struct {
	taskCh      chan func()
	closeCh     chan struct{}
	closeDoneCh chan struct{}
	rwmutex     sync.RWMutex
}

func NewTaskPool(taskNumber int) *TaskPool {
	tp := &TaskPool{
		taskCh:      make(chan func(), taskNumber),
		closeCh:     make(chan struct{}),
		closeDoneCh: make(chan struct{}),
	}

	go tp.runTask(taskNumber)
	return tp
}

func (tp *TaskPool) runTask(taskNumber int) {
	wg := sync.WaitGroup{}
	wg.Add(taskNumber)

	for i := 0; i < taskNumber; i++ {
		go func() {
			defer wg.Done()
			for task := range tp.taskCh {
				task()
			}
		}()
	}

	go func() {
		wg.Wait()
		close(tp.closeDoneCh)
	}()
}

func (tp *TaskPool) addTask(task func()) error {
	if task == nil {
		return errors.New("invalid argument")
	}

	tp.rwmutex.RLock()
	defer tp.rwmutex.RUnlock()

	select {
	case <-tp.closeCh:
		return errors.New("pool was closed")
	default:
	}

	select {
	case tp.taskCh <- task:
		return nil
	default:
		return errors.New("buffer is full")
	}
}

func (tp *TaskPool) terminate() {
	close(tp.closeCh)

	tp.rwmutex.Lock()
	close(tp.taskCh)
	tp.rwmutex.Unlock()

	<-tp.closeDoneCh
}

func main() {
	taskNumber := 3
	pool := NewTaskPool(taskNumber)

	for i := 1; i <= taskNumber; i++ {
		err := pool.addTask(func() {
			fmt.Printf("task: %d completed\n", i)
		})

		if err != nil {
			fmt.Println(err)
		}
	}

	pool.terminate()
}
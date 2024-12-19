package main

import (
	"errors"
	"sync"
)

type WorkerPool struct {
	tasksCh chan func()
	mutex   sync.RWMutex

	closed      bool
	closeDoneCh chan struct{}
}

func NewWorkerPool(workersNumber int) (*WorkerPool, error) {
	if workersNumber <= 0 {
		return nil, errors.New("incorrect workers number")
	}

	wp := &WorkerPool{
		closeDoneCh: make(chan struct{}),
		tasksCh:     make(chan func(), workersNumber),
	}

	go wp.processTasks(workersNumber)
	return wp, nil
}

func (wp *WorkerPool) processTasks(workersNumber int) {
	var wg sync.WaitGroup
	wg.Add(workersNumber)

	for i := 0; i < workersNumber; i++ {
		go func() {
			defer wg.Done()
			for task := range wp.tasksCh {
				task()
			}
		}()
	}

	wg.Wait()
	close(wp.closeDoneCh)
}

// AddTask add task to pool
func (wp *WorkerPool) AddTask(task func()) error {
	if task == nil {
		return errors.New("incorrect task")
	}

	wp.mutex.RLock()
	defer wp.mutex.RUnlock()

	if wp.closed {
		return errors.New("pool is closed")
	}

	select {
	case wp.tasksCh <- task:
		return nil
	default:
		return errors.New("pool is full")
	}
}

// Close close pool and wait all the tasks
func (wp *WorkerPool) Close() {
	if wp.closed {
		return
	}

	wp.mutex.Lock()
	wp.closed = true
	wp.mutex.Unlock()

	close(wp.tasksCh)
	<-wp.closeDoneCh
}

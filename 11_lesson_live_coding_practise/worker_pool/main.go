package main

import (
	"errors"
	"sync"
)

type WorkerPool struct {
	tasksCh     chan func()
	closeCh     chan struct{}
	closeDoneCh chan struct{}

	mutex sync.RWMutex
}

func NewWorkerPool(workersNumber int) *WorkerPool {
	wp := &WorkerPool{
		tasksCh:     make(chan func(), workersNumber),
		closeCh:     make(chan struct{}),
		closeDoneCh: make(chan struct{}),
	}

	go wp.run(workersNumber)
	return wp
}

func (wp *WorkerPool) run(workersNumber int) {
	wg := sync.WaitGroup{}
	wg.Add(workersNumber)

	for i := 0; i < workersNumber; i++ {
		go func() {
			defer wg.Done()
			for task := range wp.tasksCh {
				task()
			}
		}()
	}

	go func() {
		wg.Wait()
		close(wp.closeDoneCh)
	}()
}

// AddTask add task to pool
func (wp *WorkerPool) AddTask(task func()) error {
	if task == nil {
		return errors.New("invalid argument")
	}

	wp.mutex.RLock()
	defer wp.mutex.RUnlock()

	select {
	case <-wp.closeCh:
		return errors.New("pool was closed")
	default:
	}

	select {
	case wp.tasksCh <- task:
		return nil
	default:
		return errors.New("buffer is full")
	}
}

// Shutdown close pool and wait all the tasks
func (wp *WorkerPool) Shutdown() {
	close(wp.closeCh)

	wp.mutex.Lock()
	close(wp.tasksCh)
	wp.mutex.Unlock()

	<-wp.closeDoneCh
}

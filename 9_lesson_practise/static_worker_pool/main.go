package main

import (
	"errors"
	"sync"
)

type WorkerPool struct {
	tasksCh     chan func()
	closeCh     chan struct{}
	closeDoneCh chan struct{}
}

func NewWorkerPool(workersNumber int) *WorkerPool {
	wp := &WorkerPool{
		tasksCh:     make(chan func(), workersNumber),
		closeCh:     make(chan struct{}),
		closeDoneCh: make(chan struct{}),
	}

	wp.initializeWorkers(workersNumber)
	return wp
}

func (wp *WorkerPool) AddTask(task func()) error {
	if task == nil {
		return errors.New("invalid task")
	}

	select {
	case <-wp.closeCh:
		return errors.New("tasks queue was closed")
	default:
	}

	// TODO

	select {
	case wp.tasksCh <- task:
		return nil
	default:
		return errors.New("tasks queue is full")
	}
}

func (wp *WorkerPool) Shutdown() {
	close(wp.closeCh)
	<-wp.closeDoneCh
}

func (wp *WorkerPool) initializeWorkers(workersNumber int) {
	wg := sync.WaitGroup{}
	wg.Add(workersNumber)

	for i := 0; i < workersNumber; i++ {
		go func() {
			defer wg.Done()

			for {
				select {
				case <-wp.closeCh:
					return
				case task := <-wp.tasksCh:
					task()
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(wp.closeDoneCh)
	}()
}

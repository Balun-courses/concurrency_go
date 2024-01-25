package main

import "context"

type WorkerPool struct {
	tasks   chan func()
	workers chan struct{}
}

func NewWorkerPool(workersNumber int) *WorkerPool {
	return &WorkerPool{
		tasks:   make(chan func(), workersNumber*2),
		workers: make(chan struct{}, workersNumber),
	}
}

func (wp *WorkerPool) handle(ctx context.Context) {
	defer func() {
		<-wp.workers
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case <-ctx.Done():
			return
		case task := <-wp.tasks:
			task()
		default:
			return
		}
	}
}

func (wp *WorkerPool) Add(ctx context.Context, task func()) error {
	if len(wp.tasks) > len(wp.workers) {
		select {
		case wp.workers <- struct{}{}:
			go wp.handle(ctx)
		default:
		}
	}

	select {
	case wp.tasks <- task:
		return nil
	default:
		return errors.New("tasks queue is full")
	}
}

package main

import (
	"fmt"
	"sync"
	"time"
)

type Scheduler struct {
	tasks       map[int]*time.Timer
	taskChannel chan func()
	mutex 		sync.Mutex
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		tasks:       make(map[int]*time.Timer),
		taskChannel: make(chan func()),
	}
}

func (s *Scheduler) AddTask(taskId int, delay time.Duration, action func()) {
	if delay < 0 || action == nil {
		return
	}

	if task, exists := s.tasks[taskId]; exists {
		task.Stop()
	}

	s.tasks[taskId] = time.AfterFunc(delay, func() {
		s.taskChannel <- action
	})
}

func (s *Scheduler) CancelTask(taskId int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if task, exists := s.tasks[taskId]; exists {
		task.Stop()
		delete(s.tasks, taskId)
	}
}

func (s *Scheduler) StartScheduler() {
	go func() {
		for action := range s.taskChannel {
			action()
		}
	}()
}

func (s *Scheduler) StopScheduler() {
	for taskId, task := range s.tasks {
		task.Stop()
		delete(s.tasks, taskId)
	}
	close(s.taskChannel)
	fmt.Println("Scheduler was terminated")
}

func main() {
	scheduler := NewScheduler()
	scheduler.StartScheduler()

	scheduler.AddTask(1, 1*time.Second, func() {
		fmt.Printf("Task %d executed after %v\n",
			1, 1*time.Second)
	})

	scheduler.AddTask(2, 2*time.Second, func() {
		fmt.Printf("Task %d executed after %v\n",
			2, 2*time.Second)
	})

	scheduler.AddTask(3, 3*time.Second, func() {
		fmt.Printf("Task %d executed after %v\n",
			3, 3*time.Second)
	})

	scheduler.AddTask(4, 4*time.Second, func() {
		fmt.Printf("Task %d executed after %v\n",
			4, 4*time.Second)
	})

	scheduler.CancelTask(4)
	time.Sleep(5 * time.Second)
	scheduler.StopScheduler()
}

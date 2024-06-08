package main

import (
	"fmt"
	"sync/atomic"
	"unsafe"
)

type item struct {
	value int
	next  unsafe.Pointer
}

type Queue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func NewQueue() Queue {
	dummy := &item{}
	return Queue{
		head: unsafe.Pointer(dummy),
		tail: unsafe.Pointer(dummy),
	}
}

func (q *Queue) Push(value int) {
	// create new node
	node := &item{value: value}

	for {
		// read tail and next reference
		tail := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer(&(*item)(tail).next)

		// if tail is the same
		if tail == atomic.LoadPointer(&q.tail) {
			if next == nil {
				// no need to fix tail, try CAS
				if atomic.CompareAndSwapPointer(&(*item)(tail).next, next, unsafe.Pointer(node)) {
					// CAS was successful, try to fix tail
					atomic.CompareAndSwapPointer(&q.tail, tail, unsafe.Pointer(node))
					return
				}
			} else {
				// try to fix tail from other goroutine
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			}
		}
	}
}

func (q *Queue) Pop() int {
	for {
		head := atomic.LoadPointer(&q.head)
		tail := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer(&(*item)(head).next)

		// if head is the same
		if head == atomic.LoadPointer(&q.head) {
			// if head and tail are some node
			if head == tail {
				if next == nil {
					// queue contains only dummy node
					return -1
				} else {
					// otherwise tail must be fixed
					atomic.CompareAndSwapPointer(&q.tail, tail, next)
				}
			} else {
				// try to pop item
				value := (*item)(next).value
				if atomic.CompareAndSwapPointer(&q.head, head, next) {
					return value
				}
			}
		}
	}
}

func main() {
	queue := NewQueue()

	queue.Push(1)
	queue.Push(2)
	queue.Push(3)

	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
}

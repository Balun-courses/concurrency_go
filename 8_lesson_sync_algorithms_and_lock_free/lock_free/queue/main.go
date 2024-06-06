package main

import "fmt"

type item struct {
	value int
	next  *item
}

type Queue struct {
	head *item
	tail *item
}

func NewQueue() Queue {
	dummy := &item{}
	return Queue{
		head: dummy,
		tail: dummy,
	}
}

func (q *Queue) Push(value int) {
	q.tail.next = &item{value: value}
	q.tail = q.tail.next
}

func (q *Queue) Pop() int {
	if q.head == q.tail {
		return -1
	}

	value := q.head.next.value
	q.head = q.head.next
	return value
}

func main() {
	queue := NewQueue()

	queue.Push(10)
	queue.Push(20)
	queue.Push(30)

	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
	fmt.Println(queue.Pop())
}

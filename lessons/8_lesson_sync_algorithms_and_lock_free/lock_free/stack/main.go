package main

import "fmt"

type item struct {
	value int
	next  *item
}

type Stack struct {
	head *item
}

func NewStack() Stack {
	return Stack{}
}

func (s *Stack) Push(value int) {
	s.head = &item{value: value, next: s.head}
}

func (s *Stack) Pop() int {
	if s.head == nil {
		return -1
	}

	value := s.head.value
	s.head = s.head.next
	return value
}

func main() {
	stack := NewStack()

	stack.Push(10)
	stack.Push(20)
	stack.Push(30)

	fmt.Println(stack.Pop())
	fmt.Println(stack.Pop())
	fmt.Println(stack.Pop())
}

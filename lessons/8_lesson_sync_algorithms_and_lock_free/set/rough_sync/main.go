package main

import (
	"fmt"
	"sync"
)

type node struct {
	next  *node
	value uint
}

type Set struct {
	mutex sync.Mutex
	head  *node
}

func NewSet() *Set {
	return &Set{
		head: &node{value: 0},
	}
}

func (s *Set) Contains(value uint) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, current := s.find(value)
	return current != nil && current.value == value
}

func (s *Set) Add(value uint) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	previous, current := s.find(value)
	if current != nil && current.value == value {
		return false
	}

	newNode := &node{value: value, next: current}
	previous.next = newNode
	return true
}

func (s *Set) Remove(value uint) bool {
	if value == 0 {
		return false
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	previous, current := s.find(value)
	if current == nil || current.value != value {
		return false
	}

	previous.next = current.next
	return true
}

func (s *Set) find(value uint) (*node, *node) {
	previous := s.head
	current := s.head.next
	for current != nil && current.value < value {
		previous = current
		current = current.next
	}

	return previous, current
}

func main() {
	set := NewSet()
	set.Add(2)
	set.Add(1)
	set.Add(3)
	set.Add(2)

	fmt.Println(set.Contains(1))
	fmt.Println(set.Contains(2))
	fmt.Println(set.Contains(3))
	fmt.Println(!set.Contains(5))

	set.Remove(2)
	set.Remove(3)
	set.Remove(1)

	fmt.Println(!set.Contains(1))
	fmt.Println(!set.Contains(2))
	fmt.Println(!set.Contains(3))
	fmt.Println(!set.Contains(5))
}

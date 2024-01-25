package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

type node struct {
	sync.Mutex
	next    unsafe.Pointer
	value   uint
	dropped bool
}

func (n *node) lock() {
	if n != nil {
		n.Lock()
	}
}

func (n *node) unlock() {
	if n != nil {
		n.Unlock()
	}
}

func (n *node) getNext() *node {
	return (*node)(atomic.LoadPointer(&n.next))
}

func (n *node) setNext(next *node) {
	atomic.StorePointer(&n.next, unsafe.Pointer(next))
}

type Set struct {
	head unsafe.Pointer
}

func NewSet() *Set {
	return &Set{
		head: unsafe.Pointer(&node{value: 0}),
	}
}

func (s *Set) Contains(value uint) bool {
	for {
		var result bool
		var validated bool
		previous, current := s.find(value)
		s.withSynchronization(previous, current, func() {
			if s.validate(previous, current) {
				validated = true
				result = current != nil && !current.dropped && current.value == value
			}
		})

		if validated {
			return result
		}
	}
}

func (s *Set) Add(value uint) bool {
	for {
		var result bool
		var validated bool
		previous, current := s.find(value)
		s.withSynchronization(previous, current, func() {
			if s.validate(previous, current) {
				validated = true
				if current == nil || current.value != value {
					newNode := &node{value: value, next: unsafe.Pointer(current)}
					previous.setNext(newNode)
					result = true
				}
			}
		})

		if validated {
			return result
		}
	}
}

func (s *Set) Remove(value uint) bool {
	if value == 0 {
		return false
	}

	for {
		var result bool
		var validated bool
		previous, current := s.find(value)
		s.withSynchronization(previous, current, func() {
			if s.validate(previous, current) {
				validated = true
				if current != nil && current.value == value {
					current.dropped = true
					previous.setNext(current.getNext())
					result = true
				}
			}
		})

		if validated {
			return result
		}
	}
}

func (s *Set) find(value uint) (*node, *node) {
	previous := (*node)(s.head)
	current := previous.getNext()
	for current != nil && current.value < value {
		previous = current
		current = current.getNext()
	}

	return previous, current
}

func (s *Set) validate(previous, current *node) bool {
	return (!previous.dropped) &&
		(current == nil || !current.dropped) &&
		(previous.getNext() == current)
}

func (s *Set) withSynchronization(previous, current *node, action func()) {
	if action == nil {
		return
	}

	previous.lock()
	current.lock()
	action()
	previous.unlock()
	current.unlock()
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

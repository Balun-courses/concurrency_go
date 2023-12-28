package main

import "sync"

type Node struct {
	Value any
	next  *Node
}

type freeList struct {
	head *Node
}

func newLinkedList() freeList {
	return freeList{}
}

func (l *freeList) push(node *Node) {
	node.next = l.head
	l.head = node
}

func (l *freeList) pop() *Node {
	if l.head != nil {
		node := l.head
		l.head = l.head.next
		return node
	}

	return nil
}

type Pool struct {
	ctr func() any

	mtx  sync.Mutex
	list freeList
}

func NewPool(ctr func() any) *Pool {
	return &Pool{
		ctr: ctr,
	}
}

func (l *Pool) Get() *Node {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	node := l.list.pop()
	if node == nil {
		node = &Node{}
	}

	node.Value = l.ctr()
	return node
}

func (l *Pool) Put(node *Node) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	l.list.push(node)
}

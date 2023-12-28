package main

// go test -bench=. pool_test.go -benchmem

import (
	"sync"
	"testing"
)

type Person struct {
	name string
}

type PersonsPool struct {
	pool sync.Pool
}

func NewPersonsPool() *PersonsPool {
	return &PersonsPool{
		pool: sync.Pool{
			New: func() interface{} { return new(Person) },
		},
	}
}
func (p *PersonsPool) Get() *Person {
	return p.pool.Get().(*Person)
}

func (p *PersonsPool) Put(person *Person) {
	p.pool.Put(person)
}

func BenchmarkWithPool(b *testing.B) {
	pool := NewPersonsPool()
	for i := 0; i < b.N; i++ {
		person := pool.Get()
		person.name = "Ivan"
		pool.Put(person)
	}
}

var gPerson *Person

func BenchmarkWithoutPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		person := &Person{name: "Ivan"}
		gPerson = person
	}
}

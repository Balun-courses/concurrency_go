package main

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestCommit(t *testing.T) {
	t.Parallel()

	inMemory := NewInMemoryStorage()
	s := NewScheduler(inMemory)

	tx1 := s.StartTransaction()
	tx1.Set("key_1", "value_1")
	tx1.Set("key_2", "value_2")
	tx1.Commit()

	tx2 := s.StartTransaction()
	assert.Equal(t, "value_1", tx2.Get("key_1"))
	assert.Equal(t, "value_2", tx2.Get("key_2"))
	tx2.Commit()
}

func TestCommitWithTheSameKeys(t *testing.T) {
	t.Parallel()

	inMemory := NewInMemoryStorage()
	s := NewScheduler(inMemory)

	tx1 := s.StartTransaction()
	tx1.Set("key_1", "value_1")
	tx1.Set("key_2", "value_2")
	tx1.Set("key_2", "new_value_2")
	tx1.Commit()

	tx2 := s.StartTransaction()
	assert.Equal(t, "value_1", tx2.Get("key_1"))
	assert.Equal(t, "new_value_2", tx2.Get("key_2"))
	tx2.Commit()
}

func TestRollback(t *testing.T) {
	t.Parallel()

	inMemory := NewInMemoryStorage()
	s := NewScheduler(inMemory)

	tx1 := s.StartTransaction()
	tx1.Set("key_1", "value_1")
	tx1.Set("key_2", "value_2")
	tx1.Rollback()

	tx2 := s.StartTransaction()
	assert.Equal(t, "", tx2.Get("key_1"))
	assert.Equal(t, "", tx2.Get("key_2"))
	tx2.Commit()
}

func TestRollbackWithTheSameKeys(t *testing.T) {
	t.Parallel()

	inMemory := NewInMemoryStorage()
	s := NewScheduler(inMemory)

	tx1 := s.StartTransaction()
	tx1.Set("key_1", "value_1")
	tx1.Set("key_2", "value_2")
	tx1.Set("key_2", "new_value_2")
	tx1.Rollback()

	tx2 := s.StartTransaction()
	assert.Equal(t, "", tx2.Get("key_1"))
	assert.Equal(t, "", tx2.Get("key_2"))
	tx2.Commit()
}

func TestDirtyRead(t *testing.T) {
	t.Parallel()

	inMemory := NewInMemoryStorage()
	s := NewScheduler(inMemory)

	wg := sync.WaitGroup{}
	wg.Add(1)

	parallelTx2 := func() {
		defer wg.Done()
		tx2 := s.StartTransaction()
		assert.Equal(t, "", tx2.Get("key_1"))
		assert.Equal(t, "", tx2.Get("key_2"))
		tx2.Commit()
	}

	tx1 := s.StartTransaction()
	tx1.Set("key_1", "value_1")
	tx1.Set("key_2", "value_2")

	go parallelTx2()
	time.Sleep(time.Millisecond * 100)

	tx1.Rollback()

	wg.Wait()
}

func TestNonRepeatableRead(t *testing.T) {
	t.Parallel()

	inMemory := NewInMemoryStorage()
	s := NewScheduler(inMemory)

	wg := sync.WaitGroup{}
	wg.Add(1)

	parallelTx2 := func() {
		defer wg.Done()
		tx2 := s.StartTransaction()
		tx2.Set("key_1", "new_value_1")
		tx2.Commit()
	}

	tx1 := s.StartTransaction()
	tx1.Set("key_1", "value_1")
	assert.Equal(t, "value_1", tx1.Get("key_1"))

	go parallelTx2()
	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, "value_1", tx1.Get("key_1"))
	tx1.Commit()

	wg.Wait()
}

func TestPhantomsRead(t *testing.T) {
	t.Parallel()

	inMemory := NewInMemoryStorage()
	s := NewScheduler(inMemory)

	wg := sync.WaitGroup{}
	wg.Add(1)

	parallelTx2 := func() {
		defer wg.Done()
		tx2 := s.StartTransaction()
		tx2.Set("key_1", "new_value_1")
		tx2.Commit()
	}

	tx1 := s.StartTransaction()
	assert.Equal(t, "", tx1.Get("key_1"))

	go parallelTx2()
	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, "", tx1.Get("key_1"))
	tx1.Commit()

	wg.Wait()
}

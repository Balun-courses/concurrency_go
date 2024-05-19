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
	assert.NoError(t, tx1.Commit())

	tx2 := s.StartTransaction()
	assert.Equal(t, "value_1", tx2.Get("key_1"))
	assert.Equal(t, "value_2", tx2.Get("key_2"))
	assert.NoError(t, tx2.Commit())
}

func TestCommitWithTheSameKeys(t *testing.T) {
	t.Parallel()

	inMemory := NewInMemoryStorage()
	s := NewScheduler(inMemory)

	tx1 := s.StartTransaction()
	tx1.Set("key_1", "value_1")
	tx1.Set("key_2", "value_2")
	tx1.Set("key_2", "new_value_2")
	assert.NoError(t, tx1.Commit())

	tx2 := s.StartTransaction()
	assert.Equal(t, "value_1", tx2.Get("key_1"))
	assert.Equal(t, "new_value_2", tx2.Get("key_2"))
	assert.NoError(t, tx2.Commit())
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
	assert.NoError(t, tx2.Commit())
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
	assert.NoError(t, tx2.Commit())
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
		assert.NoError(t, tx2.Commit())
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

	parallelTx2 := func() {
		tx2 := s.StartTransaction()
		tx2.Set("key_1", "new_value_1")
		assert.NoError(t, tx2.Commit())
	}

	tx1 := s.StartTransaction()
	tx1.Set("key_1", "value_1")
	assert.Equal(t, "value_1", tx1.Get("key_1"))

	go parallelTx2()
	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, "value_1", tx1.Get("key_1"))
	assert.Error(t, tx1.Commit(), "transactions conflict")

	// repeat
	tx1 = s.StartTransaction()
	tx1.Set("key_1", "value_1")
	assert.Equal(t, "value_1", tx1.Get("key_1"))
	assert.Equal(t, "value_1", tx1.Get("key_1"))
	assert.NoError(t, tx1.Commit(), "transactions conflict")
}

func TestPhantomsRead(t *testing.T) {
	t.Parallel()

	inMemory := NewInMemoryStorage()
	s := NewScheduler(inMemory)

	parallelTx2 := func() {
		tx2 := s.StartTransaction()
		tx2.Set("key_1", "new_value_1")
		assert.NoError(t, tx2.Commit())
	}

	tx1 := s.StartTransaction()
	assert.Equal(t, "", tx1.Get("key_1"))

	go parallelTx2()
	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, "", tx1.Get("key_1"))
	assert.NoError(t, tx1.Commit())
}

func TestWriteSkew(t *testing.T) {
	t.Parallel()

	inMemory := NewInMemoryStorage()
	s := NewScheduler(inMemory)

	tx := s.StartTransaction()
	tx.Set("x", "true")
	tx.Set("y", "true")
	assert.NoError(t, tx.Commit())

	parallelTx2 := func() {
		tx2 := s.StartTransaction()
		value := tx2.Get("y")
		if value == "true" {
			time.Sleep(time.Millisecond * 200)
			tx2.Set("x", "false")
		} else {
			t.Fail()
		}

		assert.NoError(t, tx2.Commit())
	}

	tx1 := s.StartTransaction()
	value := tx1.Get("x")
	if value == "true" {
		go parallelTx2()
		time.Sleep(time.Millisecond * 100)
		tx1.Set("y", "false")
	} else {
		t.Fail()
	}

	assert.NoError(t, tx1.Commit())

	time.Sleep(500 * time.Millisecond)

	// Anomaly with Snapshot Isolation
	tx = s.StartTransaction()
	assert.Equal(t, "false", tx.Get("x"))
	assert.Equal(t, "false", tx.Get("y"))
	assert.NoError(t, tx.Commit())
}

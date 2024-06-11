package main

import (
	"errors"
)

type scheduler interface {
	get(int32, string) string
	commit(int32, map[string]string) bool
}

type Transaction struct {
	modified map[string]string
	cached   map[string]string

	scheduler scheduler

	identifier int32
	finished   bool
}

func newTransaction(scheduler scheduler, id int32) Transaction {
	return Transaction{
		modified:   make(map[string]string),
		cached:     make(map[string]string),
		scheduler:  scheduler,
		identifier: id,
	}
}

func (t *Transaction) Set(key, value string) {
	if t.finished || key == "" || value == "" {
		return
	}

	t.modified[key] = value
}

func (t *Transaction) Get(key string) string {
	if t.finished {
		return ""
	}

	if value, found := t.modified[key]; found {
		return value
	}

	if value, found := t.cached[key]; found {
		return value
	}

	value := t.scheduler.get(t.identifier, key)
	t.cached[key] = value
	return value
}

func (t *Transaction) Commit() error {
	if t.finished {
		return nil
	}

	if len(t.modified) == 0 {
		t.finished = true
		return nil
	}

	if succeed := t.scheduler.commit(t.identifier, t.modified); !succeed {
		return errors.New("transactions conflict")
	}

	return nil
}

func (t *Transaction) Rollback() {
	if t.finished {
		return
	}

	t.finished = true
}

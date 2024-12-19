package main

const (
	shared = iota + 1
	exclusive
)

type scheduler interface {
	set(int32, string, string) txOperation
	get(int32, string) (string, txOperation)

	commit([]txOperation)
	rollback([]txOperation)
}

type txOperation struct {
	lock  *txLock
	key   string
	value *string
}

type Transaction struct {
	scheduler  scheduler
	operations []txOperation
	identifier int32
	finished   bool
}

func newTransaction(scheduler scheduler, id int32) Transaction {
	return Transaction{
		scheduler:  scheduler,
		identifier: id,
	}
}

func (t *Transaction) Set(key, value string) {
	if t.finished || key == "" || value == "" {
		return
	}

	operation := t.scheduler.set(t.identifier, key, value)
	t.operations = append(t.operations, operation)
}

func (t *Transaction) Get(key string) string {
	if t.finished {
		return ""
	}

	value, operation := t.scheduler.get(t.identifier, key)
	t.operations = append(t.operations, operation)
	return value
}

func (t *Transaction) Commit() {
	if t.finished {
		return
	}

	t.scheduler.commit(t.operations)
	t.finished = true
}

func (t *Transaction) Rollback() {
	if t.finished {
		return
	}

	t.scheduler.rollback(t.operations)
	t.finished = true
}

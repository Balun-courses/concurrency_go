package wal

import (
	"context"
	"spider/internal/database/compute"
	"spider/internal/tools"
	"sync"
	"time"
)

type fsWriter interface {
	WriteBatch([]LogRecord)
}

type WAL struct {
	fsWriter     fsWriter
	flushTimeout time.Duration
	maxBatchSize int

	mutex   sync.Mutex
	batch   []LogRecord
	batches chan []LogRecord
}

func NewWAL(fsWriter fsWriter, flushTimeout time.Duration, maxBatchSize int) *WAL {
	return &WAL{
		fsWriter:     fsWriter,
		flushTimeout: flushTimeout,
		maxBatchSize: maxBatchSize,
		batches:      make(chan []LogRecord, 2), // TODO
	}
}

func (w *WAL) StartFlushing(ctx context.Context) {
	for {
		select {
		case <-ctx.Done(): // TODO
			return
		default:
		}

		tools.WithLock(&w.mutex, func() {
			if len(w.batch) != 0 {
				w.batches <- w.batch
				w.batch = nil
			}
		})

		timer := time.NewTimer(w.flushTimeout)

		select {
		case <-ctx.Done():
			return
		case batch := <-w.batches:
			w.fsWriter.WriteBatch(batch)
		case <-timer.C:
		}
	}
}

func (w *WAL) PushSET(ctx context.Context, args []string) tools.Future[error] {
	return w.push(ctx, compute.SetCommandID, args)
}

func (w *WAL) PushDEL(ctx context.Context, args []string) tools.Future[error] {
	return w.push(ctx, compute.DelCommandID, args)
}

func (w *WAL) push(ctx context.Context, commandID int, args []string) tools.Future[error] {
	txID := ctx.Value("tx").(int64)
	record := NewLogRecord(txID, commandID, args)

	tools.WithLock(&w.mutex, func() {
		w.batch = append(w.batch, record)
		if len(w.batch) == w.maxBatchSize {
			w.batches <- w.batch
			w.batch = nil
		}
	})

	return record.Result()
}

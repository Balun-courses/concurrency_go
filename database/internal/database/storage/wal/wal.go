package wal

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"spider/internal/database/compute"
	"spider/internal/tools"
)

type fsWriter interface {
	WriteBatch([]Log)
}

type fsReader interface {
	ReadLogs() ([]LogData, error)
}

type WAL struct {
	fsWriter     fsWriter
	fsReader     fsReader
	flushTimeout time.Duration
	maxBatchSize int

	mutex   sync.Mutex
	batch   []Log
	batches chan []Log

	closeCh     chan struct{}
	closeDoneCh chan struct{}

	logger *zap.Logger
}

func NewWAL(
	fsWriter fsWriter,
	fsReader fsReader,
	flushTimeout time.Duration,
	maxBatchSize int,
	logger *zap.Logger,
) *WAL {
	wal := &WAL{
		fsWriter:     fsWriter,
		fsReader:     fsReader,
		flushTimeout: flushTimeout,
		maxBatchSize: maxBatchSize,
		batches:      make(chan []Log, 1),
		closeCh:      make(chan struct{}),
		closeDoneCh:  make(chan struct{}),
		logger:       logger,
	}

	return wal
}

func (w *WAL) Start() {
	go func() {
		defer close(w.closeDoneCh)

		for {
			select {
			case <-w.closeCh:
				w.flushBatch()
				return
			case batch := <-w.batches:
				w.fsWriter.WriteBatch(batch)
			case <-time.After(w.flushTimeout):
				w.flushBatch()
			}
		}
	}()
}

func (w *WAL) Shutdown() {
	close(w.closeCh)
	<-w.closeDoneCh
}

func (w *WAL) Set(ctx context.Context, key, value string) tools.FutureError {
	return w.push(ctx, compute.SetCommandID, []string{key, value})
}

func (w *WAL) Del(ctx context.Context, key string) tools.FutureError {
	return w.push(ctx, compute.DelCommandID, []string{key})
}

func (w *WAL) flushBatch() {
	var batch []Log
	tools.WithLock(&w.mutex, func() {
		if len(w.batch) != 0 {
			batch = w.batch
			w.batch = nil
		}
	})

	if len(batch) != 0 {
		w.fsWriter.WriteBatch(batch)
	}
}

func (w *WAL) push(ctx context.Context, commandID int, args []string) tools.FutureError {
	txID := ctx.Value("tx").(int64)
	record := NewLog(txID, commandID, args)

	tools.WithLock(&w.mutex, func() {
		w.batch = append(w.batch, record)
		if len(w.batch) == w.maxBatchSize {
			w.batches <- w.batch
			w.batch = nil
		}
	})

	return record.Result()
}

func (w *WAL) Recover() ([]LogData, error) {
	logs, err := w.fsReader.ReadLogs()
	if err != nil {
		return nil, err
	}

	return logs, nil
}

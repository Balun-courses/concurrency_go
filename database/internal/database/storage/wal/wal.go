package wal

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"spider/internal/concurrency"
	"spider/internal/database/compute"
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
		logger:       logger,
	}

	return wal
}

func (w *WAL) Start(ctx context.Context) {
	ticker := time.NewTicker(w.flushTimeout)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				w.flushBatch()
				return
			default:
			}

			select {
			case <-ctx.Done():
				w.flushBatch()
				return
			case batch := <-w.batches:
				w.fsWriter.WriteBatch(batch)
			case <-ticker.C:
				w.flushBatch()
			}
		}
	}()
}

func (w *WAL) Set(ctx context.Context, key, value string) concurrency.FutureError {
	return w.push(ctx, compute.SetCommandID, []string{key, value})
}

func (w *WAL) Del(ctx context.Context, key string) concurrency.FutureError {
	return w.push(ctx, compute.DelCommandID, []string{key})
}

func (w *WAL) flushBatch() {
	var batch []Log
	concurrency.WithLock(&w.mutex, func() {
		if len(w.batch) != 0 {
			batch = w.batch
			w.batch = nil
		}
	})

	if len(batch) != 0 {
		w.fsWriter.WriteBatch(batch)
	}
}

func (w *WAL) push(ctx context.Context, commandID int, args []string) concurrency.FutureError {
	txID := ctx.Value("tx").(int64)
	record := NewLog(txID, commandID, args)

	concurrency.WithLock(&w.mutex, func() {
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

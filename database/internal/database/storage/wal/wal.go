package wal

import (
	"context"
	"errors"
	"sync"
	"time"

	"spider/internal/common"
	"spider/internal/concurrency"
	"spider/internal/database/compute"
)

type logsWriter interface {
	Write([]WriteRequest)
}

type logsReader interface {
	Read() ([]Log, error)
}

type WAL struct {
	logsWriter logsWriter
	logsReader logsReader

	flushTimeout time.Duration
	maxBatchSize int

	batches chan []WriteRequest
	mutex   sync.Mutex
	batch   []WriteRequest
}

func NewWAL(writer logsWriter, reader logsReader, flushTimeout time.Duration, maxBatchSize int) (*WAL, error) {
	if writer == nil {
		return nil, errors.New("writer is invalid")
	}
	if reader == nil {
		return nil, errors.New("reader is invalid")
	}

	return &WAL{
		logsWriter:   writer,
		logsReader:   reader,
		flushTimeout: flushTimeout,
		maxBatchSize: maxBatchSize,
		batches:      make(chan []WriteRequest, 1),
	}, nil
}

func (w *WAL) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.flushTimeout)
		defer ticker.Stop()

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
				w.logsWriter.Write(batch)
				ticker.Reset(w.flushTimeout)
			case <-ticker.C:
				w.flushBatch()
			}
		}
	}()
}

func (w *WAL) Recover() ([]Log, error) {
	// TODO: need to compact WAL segments
	return w.logsReader.Read()
}

func (w *WAL) Set(ctx context.Context, key, value string) concurrency.FutureError {
	return w.push(ctx, compute.SetCommandID, []string{key, value})
}

func (w *WAL) Del(ctx context.Context, key string) concurrency.FutureError {
	return w.push(ctx, compute.DelCommandID, []string{key})
}

func (w *WAL) push(ctx context.Context, commandID int, args []string) concurrency.FutureError {
	txID := common.GetTxIDFromContext(ctx)
	record := NewWriteRequest(txID, commandID, args)

	concurrency.WithLock(&w.mutex, func() {
		w.batch = append(w.batch, record)
		if len(w.batch) == w.maxBatchSize {
			w.batches <- w.batch
			w.batch = nil
		}
	})

	return record.FutureResponse()
}

func (w *WAL) flushBatch() {
	var batch []WriteRequest
	concurrency.WithLock(&w.mutex, func() {
		batch = w.batch
		w.batch = nil
	})

	if len(batch) != 0 {
		w.logsWriter.Write(batch)
	}
}

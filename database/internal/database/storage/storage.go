package storage

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"spider/internal/common"
	"spider/internal/concurrency"
	"spider/internal/database/compute"
	"spider/internal/database/storage/wal"
)

var (
	ErrorNotFound  = errors.New("not found")
	ErrorMutableTX = errors.New("mutable transaction on slave")
)

type Engine interface {
	Set(context.Context, string, string)
	Get(context.Context, string) (string, bool)
	Del(context.Context, string)
}

type WAL interface {
	Recover() ([]wal.Log, error)
	Set(context.Context, string, string) concurrency.FutureError
	Del(context.Context, string) concurrency.FutureError
}

type Replica interface {
	IsMaster() bool
}

type Storage struct {
	engine    Engine
	replica   Replica
	wal       WAL
	stream    <-chan []wal.Log
	generator *IDGenerator
	logger    *zap.Logger
}

func NewStorage(engine Engine, logger *zap.Logger, options ...StorageOption) (*Storage, error) {
	if engine == nil {
		return nil, errors.New("engine is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	storage := &Storage{
		engine: engine,
		logger: logger,
	}

	for _, option := range options {
		option(storage)
	}

	var lastLSN int64
	if storage.wal != nil {
		logs, err := storage.wal.Recover()
		if err != nil {
			logger.Error("failed to recover data from WAL", zap.Error(err))
		} else {
			lastLSN = storage.applyData(logs)
		}
	}

	if storage.stream != nil {
		go func() {
			for logs := range storage.stream {
				_ = storage.applyData(logs)
			}
		}()
	}

	storage.generator = NewIDGenerator(lastLSN)
	return storage, nil
}

func (s *Storage) Set(ctx context.Context, key, value string) error {
	if s.replica != nil && !s.replica.IsMaster() {
		return ErrorMutableTX
	} else if ctx.Err() != nil {
		return ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	if s.wal != nil {
		futureResponse := s.wal.Set(ctx, key, value)
		if err := futureResponse.Get(); err != nil {
			return err
		}
	}

	s.engine.Set(ctx, key, value)
	return nil
}

func (s *Storage) Del(ctx context.Context, key string) error {
	if s.replica != nil && !s.replica.IsMaster() {
		return ErrorMutableTX
	} else if ctx.Err() != nil {
		return ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	if s.wal != nil {
		futureResponse := s.wal.Del(ctx, key)
		if err := futureResponse.Get(); err != nil {
			return err
		}
	}

	s.engine.Del(ctx, key)
	return nil
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	txID := s.generator.Generate()
	ctx = common.ContextWithTxID(ctx, txID)

	value, found := s.engine.Get(ctx, key)
	if !found {
		return "", ErrorNotFound
	}

	return value, nil
}

func (s *Storage) applyData(logs []wal.Log) int64 {
	var lastLSN int64
	for _, log := range logs {
		lastLSN = max(lastLSN, log.LSN)
		ctx := common.ContextWithTxID(context.Background(), log.LSN)
		switch log.CommandID {
		case compute.SetCommandID:
			s.engine.Set(ctx, log.Arguments[0], log.Arguments[1])
		case compute.DelCommandID:
			s.engine.Del(ctx, log.Arguments[0])
		}
	}

	return lastLSN
}

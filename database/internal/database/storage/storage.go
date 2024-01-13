package storage

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"spider/internal/tools"
)

type Engine interface {
	Set(context.Context, string, string)
	Get(context.Context, string) (string, bool)
	Del(context.Context, string)
}

type WAL interface {
	Start()
	Set(context.Context, string, string) tools.FutureError
	Del(context.Context, string) tools.FutureError
	Shutdown()
}

type Replica interface {
	Start(context.Context)
	IsMaster() bool
	Shutdown()
}

type Storage struct {
	engine  Engine
	wal     WAL
	replica Replica
	logger  *zap.Logger
}

func NewStorage(
	engine Engine,
	wal WAL,
	replica Replica,
	logger *zap.Logger,
) (*Storage, error) {
	if engine == nil {
		return nil, errors.New("engine is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Storage{
		engine:  engine,
		wal:     wal,
		logger:  logger,
		replica: replica,
	}, nil
}

func (s *Storage) Start(ctx context.Context) {
	if s.wal != nil {
		if s.replica != nil {
			if s.replica.IsMaster() {
				s.wal.Start()
			}

			s.replica.Start(ctx)
		} else {
			s.wal.Start()
		}
	}
}

func (s *Storage) Shutdown() {
	if s.wal != nil {
		if s.replica != nil {
			s.replica.Shutdown()
			if s.replica.IsMaster() {
				s.wal.Shutdown()
			}
		} else {
			s.wal.Shutdown()
		}
	}
}

func (s *Storage) Set(ctx context.Context, key, value string) error {
	if s.replica != nil && !s.replica.IsMaster() {
		return errors.New("mutable transaction on slave")
	}

	if s.wal != nil {
		future := s.wal.Set(ctx, key, value)
		if err := future.Get(); err != nil {
			return err
		}
	}

	s.engine.Set(ctx, key, value)
	return nil
}

func (s *Storage) Del(ctx context.Context, key string) error {
	if s.replica != nil && !s.replica.IsMaster() {
		return errors.New("mutable transaction on slave")
	}

	if s.wal != nil {
		future := s.wal.Del(ctx, key)
		if err := future.Get(); err != nil {
			return err
		}
	}

	s.engine.Del(ctx, key)
	return nil
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	value, _ := s.engine.Get(ctx, key)
	return value, nil
}

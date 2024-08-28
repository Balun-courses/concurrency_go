package in_memory

import (
	"context"
	"errors"
	"hash/fnv"

	"go.uber.org/zap"
)

type Engine struct {
	partitions []*HashTable
	logger     *zap.Logger
}

func NewEngine(logger *zap.Logger, options ...EngineOption) (*Engine, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	engine := &Engine{
		logger: logger,
	}

	for _, option := range options {
		option(engine)
	}

	if len(engine.partitions) == 0 {
		engine.partitions = make([]*HashTable, 1)
		engine.partitions[0] = NewHashTable()
	}

	return engine, nil
}

func (e *Engine) Set(ctx context.Context, key, value string) {
	partitionIdx := 0
	if len(e.partitions) > 1 {
		partitionIdx = e.partitionIdx(key)
	}

	partition := e.partitions[partitionIdx]
	partition.Set(key, value)

	txID := ctx.Value("tx").(int64)
	e.logger.Debug("successfull set query", zap.Int64("tx", txID))
}

func (e *Engine) Get(ctx context.Context, key string) (string, bool) {
	partitionIdx := 0
	if len(e.partitions) > 1 {
		partitionIdx = e.partitionIdx(key)
	}

	partition := e.partitions[partitionIdx]
	value, found := partition.Get(key)

	txID := ctx.Value("tx").(int64)
	e.logger.Debug("successfull get query", zap.Int64("tx", txID))
	return value, found
}

func (e *Engine) Del(ctx context.Context, key string) {
	partitionIdx := 0
	if len(e.partitions) > 1 {
		partitionIdx = e.partitionIdx(key)
	}

	partition := e.partitions[partitionIdx]
	partition.Del(key)

	txID := ctx.Value("tx").(int64)
	e.logger.Debug("successfull del query", zap.Int64("tx", txID))
}

func (e *Engine) partitionIdx(key string) int {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(key))
	return int(hash.Sum32()) % len(e.partitions)
}

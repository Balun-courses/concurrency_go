package initialization

import (
	"errors"
	"time"

	"go.uber.org/zap"

	"spider/internal/configuration"
	"spider/internal/database/storage/wal"
	"spider/internal/tools"
)

const defaultFlushingBatchSize = 100
const defaultFlushingBatchTimeout = time.Millisecond * 10
const defaultMaxSegmentSize = 10 << 20
const defaultWALDataDirectory = "./data/spider/wal"

func CreateWAL(cfg *configuration.WALConfig, logger *zap.Logger) (*wal.WAL, error) {
	if logger == nil {
		return nil, errors.New("logger is incorrect")
	} else if cfg == nil {
		return nil, nil
	}

	flushingBatchSize := defaultFlushingBatchSize
	flushingBatchTimeout := defaultFlushingBatchTimeout
	maxSegmentSize := defaultMaxSegmentSize
	dataDirectory := defaultWALDataDirectory

	if cfg.FlushingBatchLength != 0 {
		flushingBatchSize = cfg.FlushingBatchLength
	}

	if cfg.FlushingBatchTimeout != 0 {
		flushingBatchTimeout = cfg.FlushingBatchTimeout
	}

	if cfg.MaxSegmentSize != "" {
		size, err := tools.ParseSize(cfg.MaxSegmentSize)
		if err != nil {
			return nil, errors.New("max segment size is incorrect")
		}

		maxSegmentSize = size
	}

	if cfg.DataDirectory != "" {
		dataDirectory = cfg.DataDirectory
	}

	fsReader := wal.NewFSReader(dataDirectory, logger)
	fsWriter := wal.NewFSWriter(dataDirectory, maxSegmentSize, logger)
	return wal.NewWAL(fsWriter, fsReader, flushingBatchTimeout, flushingBatchSize, logger), nil
}

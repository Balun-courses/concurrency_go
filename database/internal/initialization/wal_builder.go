package initialization

import (
	"errors"
	"go.uber.org/zap"
	"spider/internal/configuration"
	"spider/internal/database/storage/wal"
	"spider/internal/tools"
	"time"
)

const defaultFlushingBatchSize = 100
const defaultFlushingBatchTimeout = time.Millisecond * 10
const defaultMaxSegmentSize = 10 << 20
const defaultDataDirectory = "./data/spider/wal"

func CreateWAL(cfg *configuration.WALConfig, logger *zap.Logger) (*wal.WAL, error) {
	flushingBatchSize := defaultFlushingBatchSize
	flushingBatchTimeout := defaultFlushingBatchTimeout
	maxSegmentSize := defaultMaxSegmentSize
	dataDirectory := defaultDataDirectory

	if cfg != nil {
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
	}

	fsReader := wal.NewFSReader(dataDirectory, logger)
	fsWriter := wal.NewFSWriter(dataDirectory, maxSegmentSize, logger)
	return wal.NewWAL(fsWriter, fsReader, flushingBatchTimeout, flushingBatchSize), nil
}

package initialization

import (
	"errors"
	"github.com/inhies/go-bytesize"
	"go.uber.org/zap"
	"spider/internal/configuration"
	"spider/internal/database/storage/wal"
	"time"
)

const defaultFlushingBatchSize = 100
const defaultFlushingBatchTimeout = time.Millisecond * 10
const defaultMaxSegmentSize = 10 << 20
const defaultDataDirectory = "/data/spider/wal"

func CreateWAL(cfg *configuration.WALConfig, logger *zap.Logger) (*wal.WAL, error) {
	flushingBatchSize := defaultFlushingBatchSize
	flushingBatchTimeout := defaultFlushingBatchTimeout
	maxSegmentSize := defaultMaxSegmentSize
	dataDirectory := defaultDataDirectory

	if cfg != nil {
		if cfg.FlushingBatchSize != 0 {
			flushingBatchSize = cfg.FlushingBatchSize
		}

		if cfg.FlushingBatchTimeout != 0 {
			flushingBatchTimeout = cfg.FlushingBatchTimeout
		}

		if cfg.MaxSegmentSize != "" {
			_, err := bytesize.Parse(cfg.MaxSegmentSize)
			if err != nil {
				return nil, errors.New("max segment size is incorrect")
			}

			// maxSegmentSize = size  TODO
		}

		if cfg.DataDirectory != "" {
			dataDirectory = cfg.DataDirectory
		}
	}

	fsWriter := wal.NewFSWriter(dataDirectory, maxSegmentSize, logger)
	return wal.NewWAL(fsWriter, flushingBatchTimeout, flushingBatchSize), nil
}

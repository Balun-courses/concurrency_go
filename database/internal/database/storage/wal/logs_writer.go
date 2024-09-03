package wal

import (
	"bytes"
	"errors"

	"go.uber.org/zap"
)

type segment interface {
	Write([]byte) error
}

type LogsWriter struct {
	segment segment
	logger  *zap.Logger
}

func NewLogsWriter(segment segment, logger *zap.Logger) (*LogsWriter, error) {
	if segment == nil {
		return nil, errors.New("segment is invalid")
	}
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &LogsWriter{
		segment: segment,
		logger:  logger,
	}, nil
}

func (w *LogsWriter) Write(requests []WriteRequest) {
	var buffer bytes.Buffer
	for idx := range requests {
		log := requests[idx].Log()
		if err := log.Encode(&buffer); err != nil {
			w.logger.Warn("failed to encode logs data", zap.Error(err))
			w.acknowledgeWrite(requests, err)
			return
		}
	}

	err := w.segment.Write(buffer.Bytes())
	if err != nil {
		w.logger.Warn("failed to write logs data", zap.Error(err))
	}

	w.acknowledgeWrite(requests, err)
}

func (w *LogsWriter) acknowledgeWrite(requests []WriteRequest, err error) {
	for idx := range requests {
		requests[idx].SetResponse(err)
	}
}

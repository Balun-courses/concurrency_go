package wal

import (
	"bytes"
	"encoding/gob"
	"go.uber.org/zap"
	"os"
	"strconv"
)

type FSWriter struct {
	segment   *os.File
	directory string
	lastLSN   int64

	segmentSize    int
	maxSegmentSize int

	logger *zap.Logger
}

func NewFSWriter(directory string, maxSegmentSize int, logger *zap.Logger) *FSWriter {
	return &FSWriter{
		directory:      directory,
		maxSegmentSize: maxSegmentSize,
		logger:         logger,
	}
}

func (w *FSWriter) WriteBatch(batch []LogRecord) {
	if w.segment == nil {
		if err := w.rotateSegment(); err != nil {
			w.acknowledgeWrite(batch, err)
			return
		}
	}

	for i := 0; i < len(batch); i++ {
		if w.segmentSize > w.maxSegmentSize {
			if err := w.rotateSegment(); err != nil {
				w.acknowledgeWrite(batch[i:], err)
				return
			}
		}

		record := batch[i]
		if err := w.writeRecord(record); err != nil {
			record.SetResult(err)
		}
	}

	err := w.segment.Sync()
	if err != nil {
		w.acknowledgeWrite(batch, err)
		w.logger.Error("failed to sync segment file", zap.Error(err))
	}

	w.acknowledgeWrite(batch, err)
}

func (w *FSWriter) writeRecord(record LogRecord) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(record.data); err != nil {
		w.logger.Warn("failed to encode record", zap.Error(err))
		return err
	}

	writtenBytes, err := w.segment.Write(buffer.Bytes())
	if err != nil {
		w.logger.Warn("failed to write record", zap.Error(err))
		return err
	}

	w.segmentSize += writtenBytes
	w.lastLSN = record.LSN()
	return nil
}

func (w *FSWriter) acknowledgeWrite(batch []LogRecord, err error) {
	for _, record := range batch {
		record.SetResult(err)
	}
}

func (w *FSWriter) rotateSegment() error {
	lastLSNStr := strconv.FormatInt(w.lastLSN, 10)
	segmentName := w.directory + "/" + lastLSNStr + ".wal"

	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	segment, err := os.OpenFile(segmentName, flags, 0644)
	if err != nil {
		w.logger.Error("failed to create wal segment", zap.Error(err))
		return err
	}

	w.segment = segment
	w.segmentSize = 0
	return nil
}

package wal

import (
	"fmt"
	"go.uber.org/zap"
	"os"
)

type FSReader struct {
	logger *zap.Logger
}

func NewFSReader(logger *zap.Logger) *FSReader {
	return &FSReader{
		logger: logger,
	}
}

func (r *FSReader) Recover(directory string) ([]LogRecord, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to scan WAL directory: %w", err)
	}

	var records []LogRecord
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := directory + "/" + file.Name()
		segmentRecords, err := r.recoverWALSegment(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to recove WAL segment: %w", err)
		}

		records = append(records, segmentRecords...)
	}

	return records, nil
}

func (r *FSReader) recoverWALSegment(filename string) ([]LogRecord, error) {
	_, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// TODO
	return nil, err
}

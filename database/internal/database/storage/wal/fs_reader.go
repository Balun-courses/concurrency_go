package wal

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"go.uber.org/zap"
	"os"
	"sort"
)

type FSReader struct {
	directory string
	logger    *zap.Logger
}

func NewFSReader(directory string, logger *zap.Logger) *FSReader {
	return &FSReader{
		directory: directory,
		logger:    logger,
	}
}

func (r *FSReader) ReadLogs() ([]LogData, error) {
	files, err := os.ReadDir(r.directory)
	if err != nil {
		return nil, fmt.Errorf("failed to scan WAL directory: %w", err)
	}

	var logs []LogData
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := r.directory + "/" + file.Name()
		segmentedLogs, err := r.readSegment(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to recove WAL segment: %w", err)
		}

		logs = append(logs, segmentedLogs...)
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].LSN < logs[j].LSN
	})

	return logs, nil
}

func (r *FSReader) readSegment(filename string) ([]LogData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var logs []LogData
	buffer := bytes.NewBuffer(data)
	for buffer.Len() > 0 {
		var batch []LogData
		decoder := gob.NewDecoder(buffer)
		if err := decoder.Decode(&batch); err != nil {
			return nil, fmt.Errorf("failed to parse logs data: %w", err)
		}

		logs = append(logs, batch...)
	}

	return logs, nil
}

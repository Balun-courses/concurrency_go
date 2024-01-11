package replication

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"os"
	"spider/internal/database/storage/wal"
)

type TCPServer interface {
	HandleQueries(ctx context.Context, handler func(context.Context, []byte) []byte) error
}

type Master struct {
	server       TCPServer
	walDirectory string
	logger       *zap.Logger
}

func NewMaster(server TCPServer, walDirectory string, logger *zap.Logger) (*Master, error) {
	if server == nil {
		return nil, errors.New("server is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Master{
		server:       server,
		walDirectory: walDirectory,
		logger:       logger,
	}, nil
}

func (m *Master) HandleSynchronizations(ctx context.Context) error {
	return m.server.HandleQueries(ctx, func(_ context.Context, requestData []byte) []byte {
		var request Request
		if err := Decode(&request, requestData); err != nil {
			m.logger.Error("failed to decode replication request", zap.Error(err))
			return nil
		}

		response := m.synchronize(request)
		responseData, err := Encode(&response)
		if err != nil {
			m.logger.Error("failed to encode replication response", zap.Error(err))
		}

		return responseData
	})
}

func (m *Master) synchronize(request Request) Response {
	var response Response
	filename, err := wal.SegmentUpperBound(m.walDirectory, request.LastSegmentTimestamp)
	if err != nil {
		// TODO
		return response
	}

	if filename == "" {
		// TODO
		return response
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		// TODO
		return response
	}

	response.Succeed = true
	response.SegmentData = data
	response.SegmentTimestamp = m.extractTimestampFromSegmentFilename(filename)
	return response
}

func (m *Master) extractTimestampFromSegmentFilename(filename string) int64 {
	idx := 0
	for idx < len(filename) && filename[idx] != '_' {
		idx++
	}

	idx++
	var timestamp int64
	for idx < len(filename) && filename[idx] != '.' {
		number := filename[idx] - '0'
		timestamp = timestamp*10 + int64(number)
		idx++
	}

	return timestamp
}

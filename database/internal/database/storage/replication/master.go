package replication

import (
	"context"
	"errors"
	"fmt"
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
	segmentName, err := wal.SegmentUpperBound(m.walDirectory, request.LastSegmentName)
	if err != nil {
		m.logger.Error("failed to find WAL segment", zap.Error(err))
		return response
	}

	if segmentName == "" {
		response.Succeed = true
		return response
	}

	filename := fmt.Sprintf("%s/%s", m.walDirectory, segmentName)
	data, err := os.ReadFile(filename)
	if err != nil {
		m.logger.Error("failed to read WAL segment", zap.Error(err))
		return response
	}

	response.Succeed = true
	response.SegmentData = data
	response.SegmentName = filename
	return response
}

package replication

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.uber.org/zap"

	"spider/internal/database/storage/wal"
)

type TCPServer interface {
	HandleQueries(context.Context, func(context.Context, []byte) []byte)
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

func (m *Master) Start(ctx context.Context) {
	m.server.HandleQueries(ctx, func(ctx context.Context, requestData []byte) []byte {
		if ctx.Err() != nil {
			return nil
		}

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

func (m *Master) IsMaster() bool {
	return true
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

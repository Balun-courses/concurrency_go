package replication

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"go.uber.org/zap"
	"spider/internal/database/storage/wal"
	"time"
)

type TCPClient interface {
	Send([]byte) ([]byte, error)
}

type Slave struct {
	logger        *zap.Logger
	client        TCPClient
	stream        chan []wal.LogData
	syncInterval  time.Duration
	lastSegmentTS int64
}

func NewSlave(client TCPClient, syncInterval time.Duration, logger *zap.Logger) (*Slave, error) {
	if client == nil {
		return nil, errors.New("client is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Slave{
		client:       client,
		logger:       logger,
		stream:       make(chan []wal.LogData, 1),
		syncInterval: syncInterval,
	}, nil
}

func (s *Slave) ReplicationStream() <-chan []wal.LogData {
	return s.stream
}

func (s *Slave) StartSynchronization(ctx context.Context) {
	defer close(s.stream)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(s.syncInterval):
			s.synchronize()
		}
	}
}

func (s *Slave) synchronize() {
	request := NewRequest(s.lastSegmentTS)
	requestData, err := Encode(&request)
	if err != nil {
		s.logger.Error("failed to encode replication request", zap.Error(err))
	}

	responseData, err := s.client.Send(requestData)
	if err != nil {
		s.logger.Error("failed to send replication request", zap.Error(err))
	}

	var response Response
	if err = Decode(&response, responseData); err != nil {
		s.logger.Error("failed to decode replication response", zap.Error(err))
	}

	if !response.Succeed {
		s.logger.Error("replication error from master")
	} else if response.SegmentTimestamp != 0 {
		var logs []wal.LogData
		buffer := bytes.NewBuffer(response.SegmentData)
		decoder := gob.NewDecoder(buffer)
		if err = decoder.Decode(&logs); err != nil {
			s.logger.Error("failed to decode replicated logs", zap.Error(err))
		}

		s.lastSegmentTS = response.SegmentTimestamp
		s.stream <- logs
	}
}

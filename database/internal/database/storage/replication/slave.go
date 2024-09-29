package replication

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"spider/internal/database/filesystem"
	"spider/internal/database/storage/wal"
)

type TCPClient interface {
	Send([]byte) ([]byte, error)
	Close()
}

type Slave struct {
	client TCPClient
	stream chan []wal.Log

	syncInterval    time.Duration
	walDirectory    string
	lastSegmentName string

	logger *zap.Logger
}

func NewSlave(
	client TCPClient,
	walDirectory string,
	syncInterval time.Duration,
	logger *zap.Logger,
) (*Slave, error) {
	if client == nil {
		return nil, errors.New("client is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	segmentName, err := filesystem.SegmentLast(walDirectory)
	if err != nil {
		logger.Error("failed to find last WAL segment", zap.Error(err))
	}

	return &Slave{
		client:          client,
		stream:          make(chan []wal.Log),
		syncInterval:    syncInterval,
		walDirectory:    walDirectory,
		lastSegmentName: segmentName,
		logger:          logger,
	}, nil
}

func (s *Slave) Start(ctx context.Context) {
	ticker := time.NewTicker(s.syncInterval)
	defer func() {
		ticker.Stop()
		s.client.Close()
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.synchronize()
			}
		}
	}()
}

func (s *Slave) IsMaster() bool {
	return false
}

func (s *Slave) ReplicationStream() <-chan []wal.Log {
	return s.stream
}

func (s *Slave) synchronize() {
	request := NewRequest(s.lastSegmentName)
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

	if response.Succeed {
		s.handleResponse(response)
	} else {
		s.logger.Error("failed to apply replication data: master error")
	}
}

func (s *Slave) handleResponse(response Response) {
	if response.SegmentName == "" {
		s.logger.Debug("no changes from replication")
		return
	}

	if err := s.saveWALSegment(response.SegmentName, response.SegmentData); err != nil {
		s.logger.Error("failed to apply replication data", zap.Error(err))
	}

	if err := s.writeDataToStream(response.SegmentData); err != nil {
		s.logger.Error("failed to write data to stream", zap.Error(err))
	}

	s.lastSegmentName = response.SegmentName
}

func (s *Slave) saveWALSegment(segmentName string, segmentData []byte) error {
	filename := fmt.Sprintf("%s/%s", s.walDirectory, segmentName)
	file, err := filesystem.CreateFile(filename)
	if err != nil {
		return fmt.Errorf("failed to create wal segment: %w", err)
	}

	_, err = filesystem.WriteFile(file, segmentData)
	return err
}

func (s *Slave) writeDataToStream(segmentData []byte) error {
	var logs []wal.Log
	buffer := bytes.NewBuffer(segmentData)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&logs); err != nil {
		return fmt.Errorf("failed to decode data: %w", err)
	}

	s.stream <- logs
	return nil
}

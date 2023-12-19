package network

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"spider/internal/tools"
	"sync"
	"time"
)

type TCPHandler = func(context.Context, []byte) []byte

type TCPServer struct {
	address     string
	semaphore   tools.Semaphore
	idleTimeout time.Duration
	messageSize int
	logger      *zap.Logger
}

func NewTCPServer(
	address string,
	maxConnectionsNumber int,
	maxMessageSize int,
	idleTimeout time.Duration,
	logger *zap.Logger,
) (*TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	if maxConnectionsNumber <= 0 {
		return nil, errors.New("invalid number of max connections")
	}

	return &TCPServer{
		address:     address,
		semaphore:   tools.NewSemaphore(maxConnectionsNumber),
		idleTimeout: idleTimeout,
		messageSize: maxMessageSize,
		logger:      logger,
	}, nil
}

func (s *TCPServer) HandleQueries(ctx context.Context, handler TCPHandler) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for {
			s.semaphore.Acquire()
			connection, err := listener.Accept()
			if err != nil || ctx.Err() != nil {
				return
			} else if err != nil {
				s.logger.Warn("failed to accept", zap.Error(err))
			}

			go func() {
				wg.Add(1)
				defer func() {
					wg.Done()
					s.semaphore.Release()
				}()

				if ctx.Err() == nil {
					s.handleConnection(ctx, connection, handler)
				}
			}()
		}
	}()

	go func() {
		defer wg.Done()

		<-ctx.Done()
		if err := listener.Close(); err != nil {
			s.logger.Warn("failed to close listener", zap.Error(err))
		}
	}()

	wg.Wait()
	return nil
}

func (s *TCPServer) handleConnection(ctx context.Context, connection net.Conn, handler TCPHandler) {
	request := make([]byte, s.messageSize)

	for {
		if err := connection.SetReadDeadline(time.Now().Add(s.idleTimeout)); err != nil {
			s.logger.Warn("failed to set read deadline", zap.Error(err))
			break
		}

		count, err := connection.Read(request)
		if err != nil && err != io.EOF {
			s.logger.Warn("failed to read", zap.Error(err))
			break
		}

		response := handler(ctx, request[:count])
		if _, err := connection.Write(response); err != nil {
			s.logger.Warn("failed to write", zap.Error(err))
			break
		}
	}

	if err := connection.Close(); err != nil {
		s.logger.Warn("failed to close connection", zap.Error(err))
	}
}

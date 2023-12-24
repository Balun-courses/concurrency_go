package network

import (
	"fmt"
	"net"
	"time"
)

type TCPClient struct {
	connection     net.Conn
	maxMessageSize int
	idleTimeout    time.Duration
}

func NewTCPClient(address string, maxMessageSize int, idleTimeout time.Duration) (*TCPClient, error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return &TCPClient{
		connection:     connection,
		maxMessageSize: maxMessageSize,
		idleTimeout:    idleTimeout,
	}, nil
}

func (c *TCPClient) Send(request []byte) ([]byte, error) {
	if err := c.connection.SetDeadline(time.Now().Add(c.idleTimeout)); err != nil {
		return nil, err
	}

	if _, err := c.connection.Write(request); err != nil {
		return nil, err
	}

	response := make([]byte, c.maxMessageSize)
	count, err := c.connection.Read(response)
	if err != nil {
		return nil, err
	}

	return response[:count], nil
}

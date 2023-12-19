package network

import (
	"fmt"
	"net"
	"time"
)

type TCPClient struct {
	connection  net.Conn
	idleTimeout time.Duration
}

func NewTCPClient(address string, idleTimeout time.Duration) (*TCPClient, error) {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return &TCPClient{
		connection:  connection,
		idleTimeout: idleTimeout,
	}, nil
}

func (c *TCPClient) Send(request []byte) ([]byte, error) {
	if err := c.connection.SetReadDeadline(time.Now().Add(c.idleTimeout)); err != nil {
		return nil, err
	}

	if _, err := c.connection.Write(request); err != nil {
		return nil, err
	}

	response := make([]byte, 2048)
	count, err := c.connection.Read(response)
	if err != nil {
		return nil, err
	}

	return response[:count], nil
}

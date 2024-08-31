package network

import "time"

const defaultBufferSize = 4 << 10

type TCPClientOption func(*TCPClient)

func WithClientIdleTimeout(timeout time.Duration) TCPClientOption {
	return func(client *TCPClient) {
		client.idleTimeout = timeout
	}
}

func WithClientBufferSize(size uint) TCPClientOption {
	return func(client *TCPClient) {
		client.bufferSize = int(size)
	}
}

type TCPServerOption func(*TCPServer)

func WithServerIdleTimeout(timeout time.Duration) TCPServerOption {
	return func(server *TCPServer) {
		server.idleTimeout = timeout
	}
}

func WithServerBufferSize(size uint) TCPServerOption {
	return func(server *TCPServer) {
		server.bufferSize = int(size)
	}
}

func WithServerMaxConnectionsNumber(count uint) TCPServerOption {
	return func(server *TCPServer) {
		server.maxConnections = int(count)
	}
}

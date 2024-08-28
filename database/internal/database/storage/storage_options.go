package storage

import "spider/internal/database/storage/wal"

type StorageOption func(*Storage)

func WithReplicationStream(stream <-chan []wal.LogData) StorageOption {
	return func(storage *Storage) {
		storage.stream = stream
	}
}

func WithWAL(wal WAL) StorageOption {
	return func(storage *Storage) {
		storage.wal = wal
	}
}

func WithReplication(replica Replica) StorageOption {
	return func(storage *Storage) {
		storage.replica = replica
	}
}

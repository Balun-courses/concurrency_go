package storage

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"spider/internal/database/storage/wal"
)

func TestWithReplicationStream(t *testing.T) {
	t.Parallel()

	stream := make(<-chan []wal.Log)
	option := WithReplicationStream(stream)

	var storage Storage
	option(&storage)

	assert.Equal(t, stream, storage.stream)
}

func TestWithWAL(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	wal := NewMockWAL(ctrl)
	option := WithWAL(wal)

	var storage Storage
	option(&storage)

	assert.Equal(t, wal, storage.wal)
}

func TestWithReplication(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	replica := NewMockReplica(ctrl)
	option := WithReplication(replica)

	var storage Storage
	option(&storage)

	assert.Equal(t, replica, storage.replica)
}

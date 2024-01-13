package wal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSegmentUpperBound(t *testing.T) {
	t.Parallel()

	filename, err := SegmentUpperBound("test_data", "wal_0.log")
	require.NoError(t, err)
	require.Equal(t, "wal_1000.log", filename)

	filename, err = SegmentUpperBound("test_data", "wal_1000.log")
	require.NoError(t, err)
	require.Equal(t, "wal_2000.log", filename)

	filename, err = SegmentUpperBound("test_data", "wal_2000.log")
	require.NoError(t, err)
	require.Equal(t, "wal_3000.log", filename)

	filename, err = SegmentUpperBound("test_data", "wal_3000.log")
	require.NoError(t, err)
	require.Equal(t, "", filename)
}

func TestSegmentLast(t *testing.T) {
	t.Parallel()

	filename, err := SegmentLast("test_data")
	require.NoError(t, err)
	require.Equal(t, "wal_3000.log", filename)
}

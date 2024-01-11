package wal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSegmentUpperBound(t *testing.T) {
	t.Parallel()

	filename, err := SegmentUpperBound("test_data", 0)
	require.NoError(t, err)
	require.Equal(t, "wal_1000.log", filename)

	filename, err = SegmentUpperBound("test_data", 1000)
	require.NoError(t, err)
	require.Equal(t, "wal_2000.log", filename)

	filename, err = SegmentUpperBound("test_data", 2000)
	require.NoError(t, err)
	require.Equal(t, "wal_3000.log", filename)

	filename, err = SegmentUpperBound("test_data", 3000)
	require.NoError(t, err)
	require.Equal(t, "", filename)
}

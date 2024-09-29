package filesystem

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSegmentUpperBound(t *testing.T) {
	t.Parallel()

	filename, err := SegmentNext("test_data", "wal_0.log")
	require.NoError(t, err)
	require.Equal(t, "wal_1000.log", filename)

	filename, err = SegmentNext("test_data", "wal_1000.log")
	require.NoError(t, err)
	require.Equal(t, "wal_2000.log", filename)

	filename, err = SegmentNext("test_data", "wal_2000.log")
	require.NoError(t, err)
	require.Equal(t, "", filename)

	filename, err = SegmentNext("test_data", "wal_3000.log")
	require.NoError(t, err)
	require.Equal(t, "", filename)
}

func TestSegmentLast(t *testing.T) {
	t.Parallel()

	filename, err := SegmentLast("test_data")
	require.NoError(t, err)
	require.Equal(t, "wal_3000.log", filename)
}

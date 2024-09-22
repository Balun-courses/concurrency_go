package filesystem

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSegmentWrite(t *testing.T) {
	t.Parallel()

	const testWALDirectory = "temp_test_data"
	err := os.Mkdir(testWALDirectory, os.ModePerm)
	require.NoError(t, err)

	defer func() {
		err := os.RemoveAll(testWALDirectory)
		require.NoError(t, err)
	}()

	const maxSegmentSize = 10
	segment := NewSegment(testWALDirectory, maxSegmentSize)

	now = func() time.Time {
		return time.Unix(1, 0)
	}

	err = segment.Write([]byte("aaaaa"))
	require.NoError(t, err)
	err = segment.Write([]byte("bbbbb"))
	require.NoError(t, err)

	now = func() time.Time {
		return time.Unix(2, 0)
	}

	err = segment.Write([]byte("ccccc"))
	require.NoError(t, err)

	stat, err := os.Stat(testWALDirectory + "/wal_1000.log")
	require.NoError(t, err)
	assert.Equal(t, int64(10), stat.Size())

	stat, err = os.Stat(testWALDirectory + "/wal_2000.log")
	require.NoError(t, err)
	assert.Equal(t, int64(5), stat.Size())
}

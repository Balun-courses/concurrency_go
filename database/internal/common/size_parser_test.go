package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSizeWithBytes(t *testing.T) {
	t.Parallel()

	size, err := ParseSize("20B")
	require.NoError(t, err)
	require.Equal(t, 20, size)

	size, err = ParseSize("20b")
	require.NoError(t, err)
	require.Equal(t, 20, size)

	size, err = ParseSize("20")
	require.NoError(t, err)
	require.Equal(t, 20, size)
}

func TestParseSizeWithKiloBytes(t *testing.T) {
	t.Parallel()

	size, err := ParseSize("20KB")
	require.NoError(t, err)
	require.Equal(t, 20*1024, size)

	size, err = ParseSize("20Kb")
	require.NoError(t, err)
	require.Equal(t, 20*1024, size)

	size, err = ParseSize("20kb")
	require.NoError(t, err)
	require.Equal(t, 20*1024, size)
}

func TestParseSizeWithMegaBytes(t *testing.T) {
	t.Parallel()

	size, err := ParseSize("20MB")
	require.NoError(t, err)
	require.Equal(t, 20*1024*1024, size)

	size, err = ParseSize("20Mb")
	require.NoError(t, err)
	require.Equal(t, 20*1024*1024, size)

	size, err = ParseSize("20mb")
	require.NoError(t, err)
	require.Equal(t, 20*1024*1024, size)
}

func TestParseSizeWithGigaBytes(t *testing.T) {
	t.Parallel()

	size, err := ParseSize("20GB")
	require.NoError(t, err)
	require.Equal(t, 20*1024*1024*1024, size)

	size, err = ParseSize("20Gb")
	require.NoError(t, err)
	require.Equal(t, 20*1024*1024*1024, size)

	size, err = ParseSize("20gb")
	require.NoError(t, err)
	require.Equal(t, 20*1024*1024*1024, size)
}

func TestParseIncorrectSize(t *testing.T) {
	t.Parallel()

	_, err := ParseSize("-20")
	require.Error(t, err)

	_, err = ParseSize("-20TB")
	require.Error(t, err)

	_, err = ParseSize("abc")
	require.Error(t, err)
}

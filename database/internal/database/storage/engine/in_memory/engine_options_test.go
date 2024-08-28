package in_memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithPartitions(t *testing.T) {
	t.Parallel()

	partitionNumber := 10
	option := WithPartitions(uint(partitionNumber))

	var engine Engine
	option(&engine)

	assert.Equal(t, partitionNumber, len(engine.partitions))
}

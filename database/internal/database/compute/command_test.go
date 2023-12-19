package compute

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCommandNameToCommandID(t *testing.T) {
	t.Parallel()

	require.Equal(t, SetCommandID, CommandNameToCommandID("SET"))
	require.Equal(t, GetCommandID, CommandNameToCommandID("GET"))
	require.Equal(t, DelCommandID, CommandNameToCommandID("DEL"))
	require.Equal(t, UnknownCommandID, CommandNameToCommandID("TRUNCATE"))
}

package compute

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandNameToCommandID(t *testing.T) {
	t.Parallel()

	require.Equal(t, SetCommandID, commandNameToCommandID("SET"))
	require.Equal(t, GetCommandID, commandNameToCommandID("GET"))
	require.Equal(t, DelCommandID, commandNameToCommandID("DEL"))
	require.Equal(t, UnknownCommandID, commandNameToCommandID("TRUNCATE"))
}

func TestCommandArgumentsNumber(t *testing.T) {
	t.Parallel()

	require.Equal(t, setCommandArgumentsNumber, commandArgumentsNumber(SetCommandID))
	require.Equal(t, getCommandArgumentsNumber, commandArgumentsNumber(GetCommandID))
	require.Equal(t, delCommandArgumentsNumber, commandArgumentsNumber(DelCommandID))
	require.Equal(t, 0, commandNameToCommandID("TRUNCATE"))
}

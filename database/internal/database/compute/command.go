package compute

const (
	UnknownCommandID = iota
	SetCommandID
	GetCommandID
	DelCommandID
)

var (
	UnknownCommand = "UNKNOWN"
	SetCommand     = "SET"
	GetCommand     = "GET"
	DelCommand     = "DEL"
)

var commandNamesToId = map[string]int{
	UnknownCommand: UnknownCommandID,
	SetCommand:     SetCommandID,
	GetCommand:     GetCommandID,
	DelCommand:     DelCommandID,
}

func CommandNameToCommandID(command string) int {
	status, found := commandNamesToId[command]
	if !found {
		return UnknownCommandID
	}

	return status
}

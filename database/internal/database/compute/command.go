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

var namesToId = map[string]int{
	SetCommand: SetCommandID,
	GetCommand: GetCommandID,
	DelCommand: DelCommandID,
}

func commandNameToCommandID(command string) int {
	status, found := namesToId[command]
	if !found {
		return UnknownCommandID
	}

	return status
}

const (
	setCommandArgumentsNumber = 2
	getCommandArgumentsNumber = 1
	delCommandArgumentsNumber = 1
)

var argumentsNumber = map[int]int{
	SetCommandID: setCommandArgumentsNumber,
	GetCommandID: getCommandArgumentsNumber,
	DelCommandID: delCommandArgumentsNumber,
}

func commandArgumentsNumber(commandID int) int {
	return argumentsNumber[commandID]
}

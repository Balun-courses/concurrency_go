package compute

type Query struct {
	commandID int
	arguments []string
}

func NewQuery(commandID int, arguments []string) Query {
	return Query{
		commandID: commandID,
		arguments: arguments,
	}
}

func (c *Query) CommandID() int {
	return c.commandID
}

func (c *Query) Arguments() []string {
	return c.arguments
}

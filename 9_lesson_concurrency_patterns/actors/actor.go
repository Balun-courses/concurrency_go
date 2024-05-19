package main

const inboxSize = 10

type Executor interface {
	Execute(Message)
}

type Message struct {
	To   string
	From string
	Body string
}

type actor struct {
	address  string
	inbox    chan Message
	executor Executor
}

func newActor(address string, executor Executor) *actor {
	obj := &actor{
		address:  address,
		inbox:    make(chan Message, inboxSize),
		executor: executor,
	}

	go obj.loop()
	return obj
}

func (a *actor) loop() {
	for message := range a.inbox {
		a.executor.Execute(message)
	}
}

func (a *actor) receive(message Message) {
	a.inbox <- message
}

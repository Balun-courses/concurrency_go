package main

import (
	"fmt"
	"time"
)

const (
	NetworkActor   = "network"
	ProcessorActor = "processor"
	SenderActor    = "sender"
)

type MessageProcessor struct{}

func (p *MessageProcessor) Execute(income Message) {
	fmt.Printf("received message from [%s:%s]: %s\n", income.From, income.To, income.Body)

	outcome := Message{
		From: income.To,
		To:   SenderActor,
		Body: "processed_message",
	}

	if err := GetActorManager().SendMessage(outcome); err != nil {
		fmt.Printf("failed to send message: %s", err.Error())
	}
}

type MessageSender struct{}

func (s *MessageSender) Execute(income Message) {
	fmt.Printf("received message [%s:%s]: %s\n", income.From, income.To, income.Body)
	fmt.Printf("message successfully sent")
}

func main() {
	manager := GetActorManager()
	if err := manager.CreateActor(ProcessorActor, &MessageProcessor{}); err != nil {
		panic("failed to create actor")
	}

	if err := manager.CreateActor(SenderActor, &MessageSender{}); err != nil {
		panic("failed to create actor")
	}

	message := Message{
		From: NetworkActor,
		To:   ProcessorActor,
		Body: "received_message",
	}

	if err := manager.SendMessage(message); err != nil {
		fmt.Printf("failed to send message: %s", err.Error())
	}

	time.Sleep(time.Second)
}

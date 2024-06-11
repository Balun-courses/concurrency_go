package main

import (
	"errors"
	"sync"
)

var once sync.Once
var instance *ActorManager

type ActorManager struct {
	mutex  sync.RWMutex
	actors map[string]*actor
}

func GetActorManager() *ActorManager {
	once.Do(func() {
		instance = &ActorManager{
			actors: make(map[string]*actor),
		}
	})

	return instance
}

func (am *ActorManager) CreateActor(address string, executor Executor) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if _, found := am.actors[address]; found {
		return errors.New("already exists")
	}

	am.actors[address] = newActor(address, executor)
	return nil
}

func (am *ActorManager) SendMessage(message Message) error {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	obj, found := am.actors[message.To]
	if !found {
		return errors.New("not found")
	}

	obj.receive(message)
	return nil
}

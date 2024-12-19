package main

// Need to show solution

type Singleton struct{}

var instance *Singleton

func GetInstance() *Singleton {
	if instance == nil {
		instance = &Singleton{}
	}

	return instance
}

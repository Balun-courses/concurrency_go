package main

import (
	"fmt"
	"reflect"
)

func main() {
	ch := make(chan int, 1)
	vch := reflect.ValueOf(ch)

	succeed := vch.TrySend(reflect.ValueOf(100))
	fmt.Println(succeed, vch.Len(), vch.Cap())

	branches := []reflect.SelectCase{
		{Dir: reflect.SelectDefault},
		{Dir: reflect.SelectRecv, Chan: vch},
	}

	index, vRecv, recvOk := reflect.Select(branches)
	fmt.Println(index, vRecv, recvOk)

	vch.Close()
}

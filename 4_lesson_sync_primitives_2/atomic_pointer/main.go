package main

import (
	"sync/atomic"
	"unsafe"
)

func main() {
	{
		var value1 int32 = 100
		var value2 int32 = 100
		pointer := unsafe.Pointer(&value1)
		atomic.StorePointer(&pointer, unsafe.Pointer(&value2))
	}
	{
		var value1 int32 = 100
		var value2 int32 = 100
		var pointer atomic.Pointer[int32]
		pointer.Store(&value1)
		pointer.Store(&value2)
	}
}

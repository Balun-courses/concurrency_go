package order_test

import (
	"fmt"
	"testing"
)

// go test -v

func trace(name string) func() {
	fmt.Printf("%s entered\n", name)
	return func() {
		fmt.Printf("%s returned\n", name)
	}
}

func Test1(t *testing.T) {
	defer trace("Test1")()
}

func Test2(t *testing.T) {
	defer trace("Test2")()
	t.Parallel()
}

func Test3(t *testing.T) {
	defer trace("Test3")()
}

func Test5(t *testing.T) {
	defer trace("Test5")()
	t.Parallel()
}

func Test4(t *testing.T) {
	defer trace("Test4")()
}

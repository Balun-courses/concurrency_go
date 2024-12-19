package second_package

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// go test -p 4 -v ./...

//    -p n
//        the number of programs, such as build commands or
//        test binaries, that can be run in parallel.
//        The default is the number of CPUs available.

//  What happens if we specify -p=1?
//  There would be only one process running tests, so all tests would be run sequentially, one package at a time.

func Test1(t *testing.T) {
	assert.Equal(t, 10, 5+5)
}

func Test2(t *testing.T) {
	assert.Equal(t, 10, 5+5)
}

func Test5(t *testing.T) {
	assert.Equal(t, 10, 5+5)
}

func Test4(t *testing.T) {
	assert.Equal(t, 10, 5+5)
}

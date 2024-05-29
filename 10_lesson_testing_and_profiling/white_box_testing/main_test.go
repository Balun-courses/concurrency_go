package white_box_testing

import (
	"testing"
	"time"
)

func ProcessNumber(input int, output chan<- int) {
	go func() {
		time.Sleep(2 * time.Second)
		output <- input * 2
	}()
}

func TestProcessNumber(t *testing.T) {
	input := 5
	expected := 10
	output := make(chan int)

	ProcessNumber(input, output)

	select {
	case result := <-output:
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Test timed out")
	}
}

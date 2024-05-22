package main

func tryToReadFromChannel(ch chan string) (string, bool) {
	select {
	case value := <-ch:
		return value, true
	default:
		return "", false
	}
}

func tryToWriteFromChannel(ch chan string, value string) bool {
	select {
	case ch <- value:
		return true
	default:
		return false
	}
}

func tryToReadOrWrite(ch1 chan string, ch2 chan string) {
	select {
	case <-ch1:
	case ch2 <- "test":
	}
}

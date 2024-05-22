package main

func tryToReadFromChannel(ch chan string) (string, bool) {
	if len(ch) != 0 {
		value := <-ch
		return value, true
	} else {
		return "", false
	}
}

func tryToWriteFromChannel(ch chan string, value string) bool {
	if len(ch) < cap(ch) {
		ch <- value
		return true
	} else {
		return false
	}
}

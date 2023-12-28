package main

import "sync"

func setter(data map[string]string) {
	data["test"] = "test"
}

func main() {
	data := make(map[string]string)

	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			setter(data)
		}()
	}

	wg.Wait()
}

package main

import "sync"

type Data struct {
	sync.Mutex
	values []int
}

func (d *Data) Add(value int) {
	d.Lock()
	defer d.Unlock()

	d.values = append(d.values, value)
}

func main() {
	data := Data{}
	data.Add(100)

	data.Unlock() // Possible problem!
}

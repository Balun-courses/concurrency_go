package main

import (
	"os"
	"runtime/pprof"
)

func main() {
	fd, _ := os.Create("./cpu_profile.out")
	_ = pprof.StartCPUProfile(fd)

	// ...

	pprof.StopCPUProfile()
	_ = fd.Close()
}

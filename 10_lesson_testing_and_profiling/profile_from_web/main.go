package main

import (
	"fmt"
	"hash/fnv"
	"net/http"
	_ "net/http/pprof" // important
	"runtime"
	"sync"
)

var mutex sync.Mutex

// ab -n 100000 -c100 http://localhost:6060/headers

// go tool pprof "http://localhost:6060/debug/pprof/profile?seconds=5"
// go tool pprof "http://localhost:6060/debug/pprof/mutex?seconds=5"
// go tool pprof "http://localhost:6060/debug/pprof/block?seconds=5"

// http://localhost:6060/debug/pprof
// curl "http://localhost:6060/debug/pprof/profile?seconds=5" > profile.out

func Handle(w http.ResponseWriter, req *http.Request) {
	mutex.Lock()

	var hashes []uint32
	for i := 0; i < 1000000; i++ {
		hash := fnv.New32a()
		_, _ = hash.Write([]byte("value"))
		hashes = append(hashes, hash.Sum32())
	}

	mutex.Unlock()

	_, _ = fmt.Fprintf(w, "response\n")
}

func main() {
	runtime.SetBlockProfileRate(1)     // by default turned off
	runtime.SetMutexProfileFraction(1) // by default turned off

	http.HandleFunc("/headers", Handle)

	/*
		http.HandleFunc("/debug/pprof/", pprof.Index)
		http.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		http.HandleFunc("/debug/pprof/profile", pprof.Profile)
		http.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		http.HandleFunc("/debug/pprof/trace", pprof.Trace)
	*/

	_ = http.ListenAndServe("localhost:6060", nil)
}

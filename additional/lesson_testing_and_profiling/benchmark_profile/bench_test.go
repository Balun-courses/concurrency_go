package perftest

import (
	"testing"
)

// go test -bench=.
// go test -bench=. -benchmem
// go test -bench=. -cpuprofile=profile.out
// go tool pprof profile.out
//		top
//      top 20
//      top 20 -cum
//      list cp
//      list runtime.mallocgc
//		web
//		trim=false
//      web
// go tool pprof -http :8081 profile.out

// go test -bench=. -memprofile=profile.out
// go tool pprof profile.out

// go test -bench=. -benchmem --count=6 > profile.out
// go test -bench=. -benchmem --count=6 > profile_new.out
// benchstat profile_new.out profile.out

func cp(input []string) []string {
	var output []string
	for _, value := range input {
		output = append(output, value)
	}

	return output
}

func cpOptimized(input []string) []string {
	output := make([]string, 0, len(input))
	for _, value := range input {
		output = append(output, value)
	}

	return output
}

func BenchmarkCopy(b *testing.B) {
	data := []string{"1", "2", "3", "4", "5"}
	for i := 0; i < b.N; i++ {
		cpOptimized(data)
	}
}

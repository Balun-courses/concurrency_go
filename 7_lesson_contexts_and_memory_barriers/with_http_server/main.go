package main

import (
	"context"
	"fmt"
	"net/http"
)

func main() {
	helloWorldHandler := http.HandlerFunc(handle)
	http.Handle("/welcome", injectTraceID(helloWorldHandler))
	_ = http.ListenAndServe(":8080", nil)
}

func handle(_ http.ResponseWriter, r *http.Request) {
	value, ok := r.Context().Value("trace_id").(string)
	if ok {
		fmt.Println(value)
	}

	makeRequest(r.Context())
}

func makeRequest(_ context.Context) {
	// requesting to database with context
}

func injectTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "trace_id", "12-21-33")
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

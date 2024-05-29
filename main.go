package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	rl := NewRateLimiter()
	r := mux.NewRouter()
	r.Use(RateLimitMiddleware(rl))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	http.ListenAndServe(":8080", r)
}

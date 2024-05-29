package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	persistence := NewRedisPersistence(redisAddr)
	rl := NewRateLimiter(persistence)
	r := mux.NewRouter()
	r.Use(RateLimitMiddleware(rl))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	http.ListenAndServe(":8080", r)
}
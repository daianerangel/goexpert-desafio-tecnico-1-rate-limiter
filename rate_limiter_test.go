package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter()
	handler := RateLimitMiddleware(rl)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "default")

	for i := 0; i < 10; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if i < 5 && rr.Code != http.StatusOK {
			t.Errorf("expected status OK, got %v", rr.Code)
		} else if i >= 5 && rr.Code != http.StatusTooManyRequests {
			t.Errorf("expected status Too Many Requests, got %v", rr.Code)
		}
	}
}

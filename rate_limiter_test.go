package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// BuildRedisImage builds a Docker image for Redis
func BuildRedisImage() error {
	cmd := exec.Command("docker", "build", "-t", "custom-redis:latest", ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build Redis image: %v, output: %s", err, out)
	}
	return nil
}

func StartRedis() (string, error) {
	err := BuildRedisImage()
	if err != nil {
		return "", err
	}
	cmd := exec.Command("docker", "run", "-d", "--name", "redis-test", "-p", "6379:6379", "redis:alpine")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to start Redis container: %v, output: %s", err, out)
	}
	containerID := strings.Split(string(out), "\n")[0]
	return containerID, nil
}

func StopRedis(containerID string) error {
	cmd := exec.Command("docker", "stop", containerID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop Redis container: %v, output: %s", err, out)
	}
	return nil
}

func RemoveRedis(containerID string) error {
	cmd := exec.Command("docker", "rm", containerID, "-f")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove Redis container: %v, output: %s", err, out)
	}
	return nil
}
func TestRateLimiter(t *testing.T) {
	containerID, err := StartRedis()
	if err != nil {
		t.Fatalf("could not start Redis container: %v", err)
	}
	defer StopRedis(containerID) 
	defer RemoveRedis(containerID) 

	// Wait for Redis to be ready
	time.Sleep(5 * time.Second)

	redisAddr := "localhost:6379"
	persistence := NewRedisPersistence(redisAddr)
	rl := NewRateLimiter(persistence)
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

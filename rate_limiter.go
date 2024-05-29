package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var ctx = context.Background()

type RateLimiter struct {
	persistence Persistence
	ipLimit     int
	tokenLimit  map[string]int
	blockTime   time.Duration
}

func NewRateLimiter(persistence Persistence) *RateLimiter {
	_ = godotenv.Load()
	ipLimit, _ := strconv.Atoi(os.Getenv("IP_LIMIT"))
	blockTime, _ := strconv.Atoi(os.Getenv("BLOCK_TIME"))
	tokenLimit := make(map[string]int)
	tokenLimit["default"] = 5 // Default token limit
	// Add more token limits as needed

	return &RateLimiter{
		persistence: persistence,
		ipLimit:     ipLimit,
		tokenLimit:  tokenLimit,
		blockTime:   time.Duration(blockTime) * time.Second,
	}
}

func (rl *RateLimiter) AllowRequest(ip, token string) bool {
	// Check token limit first
	if limit, exists := rl.tokenLimit[token]; exists {
		return rl.checkLimit(token, limit)
	}
	// Fallback to IP limit
	return rl.checkLimit(ip, rl.ipLimit)
}

func (rl *RateLimiter) checkLimit(key string, limit int) bool {
	count, _ := rl.persistence.Get(ctx, key)

	if count >= limit {
		return false
	}
	if err := rl.persistence.Incr(ctx, key); err != nil {
		fmt.Println("Error incrementing key in persistence:", err)
		return false
	}
	if err := rl.persistence.Expire(ctx, key, rl.blockTime); err != nil {
		fmt.Println("Error setting expiration in persistence:", err)
		return false
	}
	return true
}

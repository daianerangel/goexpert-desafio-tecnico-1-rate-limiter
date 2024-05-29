package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

type RateLimiter struct {
	redisClient *redis.Client
	ipLimit     int
	tokenLimit  map[string]int
	blockTime   time.Duration
}

func NewRateLimiter() *RateLimiter {
	_ = godotenv.Load()
	redisAddr := os.Getenv("REDIS_ADDR")
	ipLimit, _ := strconv.Atoi(os.Getenv("IP_LIMIT"))
	blockTime, _ := strconv.Atoi(os.Getenv("BLOCK_TIME"))
	tokenLimit := make(map[string]int)
	tokenLimit["default"] = 10 // Default token limit
	// Add more token limits as needed

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return &RateLimiter{
		redisClient: rdb,
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
	count, err := rl.redisClient.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		fmt.Println("Error fetching from Redis:", err)
		return false
	}
	if count >= limit {
		return false
	}
	pipe := rl.redisClient.TxPipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, rl.blockTime)
	_, err = pipe.Exec(ctx)
	if err != nil {
		fmt.Println("Error updating Redis:", err)
		return false
	}
	return true
}

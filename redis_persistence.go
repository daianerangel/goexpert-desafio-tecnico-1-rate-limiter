package main

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisPersistence struct {
	client *redis.Client
}

func NewRedisPersistence(addr string) *RedisPersistence {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisPersistence{client: client}
}

func (r *RedisPersistence) Get(ctx context.Context, key string) (int, error) {
	return r.client.Get(ctx, key).Int()
}

func (r *RedisPersistence) Incr(ctx context.Context, key string) error {
	return r.client.Incr(ctx, key).Err()
}

func (r *RedisPersistence) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

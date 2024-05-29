package main

import (
	"context"
	"time"
)

// Persistence defines the interface for the persistence layer
type Persistence interface {
	Get(ctx context.Context, key string) (int, error)
	Incr(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

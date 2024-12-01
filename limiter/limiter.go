package limiter

import (
	"context"
	"time"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string, limit int) (bool, time.Duration)
}

func NewRateLimiter(strategy RateLimiter) RateLimiter {
	return strategy
}

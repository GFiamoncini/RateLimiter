package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisLimiter struct {
	client *redis.Client
}

func NewRedisLimiter(addr, password string) *RedisLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &RedisLimiter{client: rdb}
}

func (r *RedisLimiter) Allow(ctx context.Context, key string, limit int) (bool, time.Duration) {
	pipe := r.client.TxPipeline()

	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Second)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0
	}

	if incr.Val() > int64(limit) {
		ttl, _ := r.client.TTL(ctx, key).Result()
		return false, ttl
	}

	return true, 0
}

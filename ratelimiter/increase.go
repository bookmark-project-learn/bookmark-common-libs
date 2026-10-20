package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func (r *rateLimiter) IncreaseRateLimit(ctx context.Context, key string, exp time.Duration) error {
	_, err := r.redisClient.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.Incr(ctx, key)
		p.ExpireNX(ctx, key, exp) // set the expiration time when not exists
		return nil
	})

	return err
}

package ratelimiter

import "context"

func (r *rateLimiter) GetCurrentRateLimit(ctx context.Context, key string) (int, error) {
	return r.redisClient.Get(ctx, key).Int()
}

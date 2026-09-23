package ratelimiter

import (
	"context"
	"testing"
	"time"

	redisMocks "github.com/bookmark-common-libs/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_Increase(t *testing.T) {

	testCase := []struct {
		name       string
		key        string
		exp        time.Duration
		redisMocks func(t *testing.T, ctx context.Context, key string, exp time.Duration) *redis.Client
		verifyFunc func(t *testing.T, ctx context.Context, r *redis.Client, key string, err error)
	}{
		{
			name: "set new rate and new expiration time",
			redisMocks: func(t *testing.T, ctx context.Context, key string, exp time.Duration) *redis.Client {
				mocks := redisMocks.InitMockRedis(t)
				return mocks
			},
			exp: time.Minute,
			key: "test",
			verifyFunc: func(t *testing.T, ctx context.Context, r *redis.Client, key string, err error) {
				count, err := r.Get(ctx, key).Int()
				assert.NoError(t, err)
				assert.True(t, count == 1)

				ttl, err := r.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, ttl > 0)
			},
		},
		{
			name: "increase rate limit but not change expiration time",
			redisMocks: func(t *testing.T, ctx context.Context, key string, exp time.Duration) *redis.Client {
				mocks := redisMocks.InitMockRedis(t)
				mocks.Set(ctx, key, 5, exp)
				return mocks
			},
			exp: time.Minute,
			key: "test",
			verifyFunc: func(t *testing.T, ctx context.Context, r *redis.Client, key string, err error) {
				count, err := r.Get(ctx, key).Int()
				assert.NoError(t, err)
				assert.True(t, count == 6)

				ttl, err := r.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, ttl > 0)
			},
		},
		{
			name: "falied to increase rate limit",
			redisMocks: func(t *testing.T, ctx context.Context, key string, exp time.Duration) *redis.Client {
				mocks := redisMocks.InitMockRedis(t)
				mocks.Close()
				return mocks
			},
			exp: time.Minute,
			key: "test",
			verifyFunc: func(t *testing.T, ctx context.Context, r *redis.Client, key string, err error) {
				assert.Equal(t, err, redis.ErrClosed)
			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			mocks := tc.redisMocks(t, ctx, tc.key, tc.exp)
			rateLimiter := NewRateLimiter(mocks)
			err := rateLimiter.IncreaseRateLimit(ctx, tc.key, tc.exp)
			tc.verifyFunc(t, ctx, mocks, tc.key, err)
		})
	}

}

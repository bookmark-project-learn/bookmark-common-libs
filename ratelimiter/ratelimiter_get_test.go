package ratelimiter

import (
	"context"
	"testing"
	"time"

	redisMocks "github.com/bookmark-common-libs/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestGetRateValue(t *testing.T) {
	testCase := []struct {
		name       string
		key        string
		exp        time.Duration
		redisMocks func(t *testing.T, ctx context.Context, key string, exp time.Duration) *redis.Client
		verifyFunc func(t *testing.T, ctx context.Context, r *redis.Client, result int, err error)
	}{
		{
			name: "get existing rate limit",
			redisMocks: func(t *testing.T, ctx context.Context, key string, exp time.Duration) *redis.Client {
				mocks := redisMocks.InitMockRedis(t)
				mocks.Set(ctx, key, 5, exp)
				return mocks
			},
			exp: time.Minute,
			key: "test",
			verifyFunc: func(t *testing.T, ctx context.Context, r *redis.Client, result int, err error) {
				assert.True(t, result == 5)
				assert.NoError(t, err)
			},
		},
		{
			name: "get non existing rate limit",
			redisMocks: func(t *testing.T, ctx context.Context, key string, exp time.Duration) *redis.Client {
				mocks := redisMocks.InitMockRedis(t)
				return mocks
			},
			exp: time.Minute,
			key: "test",
			verifyFunc: func(t *testing.T, ctx context.Context, r *redis.Client, result int, err error) {
				assert.Equal(t, err, redis.Nil)
			},
		},
		{
			name: "falied to get existing rate limit due to redis",
			redisMocks: func(t *testing.T, ctx context.Context, key string, exp time.Duration) *redis.Client {
				mocks := redisMocks.InitMockRedis(t)
				mocks.Close()
				return mocks
			},
			exp: time.Minute,
			key: "test",
			verifyFunc: func(t *testing.T, ctx context.Context, r *redis.Client, result int, err error) {
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
			value, err := rateLimiter.GetCurrentRateLimit(ctx, tc.key)
			tc.verifyFunc(t, ctx, mocks, value, err)
		})
	}
}

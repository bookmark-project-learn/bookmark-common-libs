package redis

import "github.com/redis/go-redis/v9"

// NewRedisClient creates a new redis client
func NewRedisClient() (*redis.Client, error) {
	cfg, err := newConfig("")
	if err != nil {
		return nil, err
	}
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	}), nil
}

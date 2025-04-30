package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheService interface {
	Get(ctx context.Context, key string) (string, bool)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

type redisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis-based cache service.
func NewRedisCache(client *redis.Client) CacheService {
	return &redisCache{client: client}
}

func (r *redisCache) Get(ctx context.Context, key string) (string, bool) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false
	} else if err != nil {
		return "", false
	}
	return val, true
}

func (r *redisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

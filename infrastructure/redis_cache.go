package infrastructure

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type RedisCache struct {
	cache *cache.Cache
}

func NewRedisCache(redisClient *redis.Client) *RedisCache {
	return &RedisCache{
		cache: cache.New(&cache.Options{
			Redis:      redisClient,
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
			Marshal:    json.Marshal,
			Unmarshal:  json.Unmarshal,
		}),
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	var value string
	err := r.cache.Get(ctx, key, &value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.cache.Delete(ctx, key)
}

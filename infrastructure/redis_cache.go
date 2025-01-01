package infrastructure

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
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

func (r *RedisCache) Get(key string) (interface{}, error) {
	var value interface{}
	err := r.cache.Get(context.Background(), key, &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (r *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	return r.cache.Set(&cache.Item{
		Ctx:   context.Background(),
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}

func (r *RedisCache) Delete(key string) error {
	return r.cache.Delete(context.Background(), key)
}

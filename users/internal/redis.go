package internal

import (
	"context"
	"time"
	"users/internal/config"

	redis "github.com/go-redis/redis/v8"
)

type RedisClient interface {
	Ping(ctx context.Context) *redis.StatusCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Close() error
}

func NewRedisClient(client *redis.Client) RedisClient {
	return client
}

// Initialize Redis connection
func InitRedis() RedisClient {
	cfg := config.LoadConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	return NewRedisClient(rdb)
}

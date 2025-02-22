package internal

import (
	"users/internal/config"

	redis "github.com/go-redis/redis/v8"
)

// Initialize Redis connection
func InitRedis() *redis.Client {
	cfg := config.LoadConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword, // Set if needed
		DB:       cfg.RedisDB,       // Default DB
	})
	return rdb
}

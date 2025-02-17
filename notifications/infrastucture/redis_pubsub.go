package infrastructure

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisEventListener struct {
	client *redis.Client
}

func NewRedisEventListener(redisAddr string) *RedisEventListener {
	client := redis.NewClient(&redis.Options{Addr: redisAddr})
	return &RedisEventListener{client: client}
}

func (r *RedisEventListener) Subscribe(channel string, handler func(string)) error {
	pubsub := r.client.Subscribe(context.Background(), channel)
	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			handler(msg.Payload)
		}
	}()
	return nil
}

func (r *RedisEventListener) Publish(channel string, message string) error {
	return r.client.Publish(context.Background(), channel, message).Err()
}

func (r *RedisEventListener) Ping(ctx context.Context) *redis.StatusCmd {
	return r.Ping(ctx)
}

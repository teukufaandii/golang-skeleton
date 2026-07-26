package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Client *redis.Client
}

// NewRedisService creates a new instance of RedisService wrapping the redis client.
func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{
		Client: client,
	}
}

// Publish sends a message to the specified channel.
func (s *RedisService) Publish(ctx context.Context, channel string, message interface{}) error {
	err := s.Client.Publish(ctx, channel, message).Err()
	if err != nil {
		return fmt.Errorf("failed to publish message to channel %s: %w", channel, err)
	}
	return nil
}

// Subscribe listens for messages on specified channels and processes them with the provided handler function.
func (s *RedisService) Subscribe(ctx context.Context, channels []string, handler func(channel string, payload string)) error {
	pubsub := s.Client.Subscribe(ctx, channels...)
	defer pubsub.Close()

	// Receive confirmation that subscription was created
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe to channels %v: %w", channels, err)
	}

	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			log.Println("Unsubscribing from Redis channels due to context cancellation")
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			handler(msg.Channel, msg.Payload)
		}
	}
}

// Set stores a key-value pair with an expiration duration.
func (s *RedisService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.Client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a string value by key.
func (s *RedisService) Get(ctx context.Context, key string) (string, error) {
	return s.Client.Get(ctx, key).Result()
}

// Del deletes specified keys from Redis.
func (s *RedisService) Del(ctx context.Context, keys ...string) error {
	return s.Client.Del(ctx, keys...).Err()
}

package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channels []string, handler func(channel string, payload string)) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
}

type redisRepository struct {
	client *redis.Client
}

// NewRedisRepository creates a new instance of RedisRepository wrapping the redis client.
func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{
		client: client,
	}
}

// Publish sends a message to the specified channel.
func (r *redisRepository) Publish(ctx context.Context, channel string, message interface{}) error {
	err := r.client.Publish(ctx, channel, message).Err()
	if err != nil {
		return fmt.Errorf("failed to publish message to channel %s: %w", channel, err)
	}
	return nil
}

// Subscribe listens for messages on specified channels and processes them with the provided handler function.
func (r *redisRepository) Subscribe(ctx context.Context, channels []string, handler func(channel string, payload string)) error {
	pubsub := r.client.Subscribe(ctx, channels...)
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
			go func(channel, payload string) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("recovered from panic in redis subscribe handler: %v", r)
					}
				}()
				handler(channel, payload)
			}(msg.Channel, msg.Payload)
		}
	}
}

// Set stores a key-value pair with an expiration duration.
func (r *redisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a string value by key.
func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del deletes specified keys from Redis.
func (r *redisRepository) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

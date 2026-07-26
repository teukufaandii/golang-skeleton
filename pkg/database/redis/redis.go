package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

// NewRedisClient initializes and returns a new Redis client instance after testing the connection with Ping.
func NewRedisClient(cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Println("Connected to Redis successfully")
	return client, nil
}

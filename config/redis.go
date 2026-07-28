package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context) (*redis.Client, error) {
	database, err := strconv.Atoi(envOrDefault("REDIS_DB", "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:         envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           database,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	return client, nil
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

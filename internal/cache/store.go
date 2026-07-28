package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store interface {
	Get(ctx context.Context, key string, destination interface{}) (bool, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	DeleteByPrefix(ctx context.Context, prefix string) error
}

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) Store {
	return &RedisStore{client: client}
}

func (s *RedisStore) Get(ctx context.Context, key string, destination interface{}) (bool, error) {
	value, err := s.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal(value, destination); err != nil {
		return false, err
	}

	return true, nil
}

func (s *RedisStore) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, key, encoded, ttl).Err()
}

func (s *RedisStore) DeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64

	for {
		keys, nextCursor, err := s.client.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := s.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			return nil
		}
	}
}

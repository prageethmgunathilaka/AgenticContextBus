package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStore wraps Redis connection
type RedisStore struct {
	client *redis.Client
}

// NewRedisStore creates a new Redis store
func NewRedisStore(addr, password string) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStore{client: client}, nil
}

// Close closes the Redis connection
func (s *RedisStore) Close() error {
	return s.client.Close()
}

// Health checks if Redis is healthy
func (s *RedisStore) Health(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

// Get retrieves a value
func (s *RedisStore) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

// Set sets a value with optional expiration
func (s *RedisStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.client.Set(ctx, key, value, expiration).Err()
}

// Delete deletes a key
func (s *RedisStore) Delete(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}

// Exists checks if a key exists
func (s *RedisStore) Exists(ctx context.Context, key string) (bool, error) {
	count, err := s.client.Exists(ctx, key).Result()
	return count > 0, err
}

// SetNX sets a key only if it doesn't exist (for idempotency)
func (s *RedisStore) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return s.client.SetNX(ctx, key, value, expiration).Result()
}

// Increment increments a counter
func (s *RedisStore) Increment(ctx context.Context, key string) (int64, error) {
	return s.client.Incr(ctx, key).Result()
}

// GetClient returns the underlying Redis client for advanced operations
func (s *RedisStore) GetClient() *redis.Client {
	return s.client
}

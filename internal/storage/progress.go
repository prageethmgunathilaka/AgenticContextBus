package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// ProgressStore interface for stream progress tracking
type ProgressStore interface {
	Get(ctx context.Context, streamID string) (*StreamProgress, error)
	Set(ctx context.Context, progress *StreamProgress, ttl time.Duration) error
	Delete(ctx context.Context, streamID string) error
}

// StreamProgress represents stream progress
type StreamProgress struct {
	StreamID      string
	Status        string
	BytesReceived int64
	TotalBytes    int64
	Progress      float64
	Checksum      string
	StartedAt     time.Time
	UpdatedAt     time.Time
}

// RedisProgressStore implements ProgressStore using Redis
type RedisProgressStore struct {
	redis *RedisStore
}

// NewRedisProgressStore creates a new Redis progress store
func NewRedisProgressStore(redis *RedisStore) *RedisProgressStore {
	return &RedisProgressStore{redis: redis}
}

// Get retrieves stream progress
func (r *RedisProgressStore) Get(ctx context.Context, streamID string) (*StreamProgress, error) {
	key := fmt.Sprintf("stream:%s", streamID)
	value, err := r.redis.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	var progress StreamProgress
	if err := json.Unmarshal([]byte(value), &progress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal progress: %w", err)
	}

	return &progress, nil
}

// Set stores stream progress
func (r *RedisProgressStore) Set(ctx context.Context, progress *StreamProgress, ttl time.Duration) error {
	key := fmt.Sprintf("stream:%s", progress.StreamID)
	data, err := json.Marshal(progress)
	if err != nil {
		return fmt.Errorf("failed to marshal progress: %w", err)
	}

	return r.redis.Set(ctx, key, string(data), ttl)
}

// Delete removes stream progress
func (r *RedisProgressStore) Delete(ctx context.Context, streamID string) error {
	key := fmt.Sprintf("stream:%s", streamID)
	return r.redis.Delete(ctx, key)
}

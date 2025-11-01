package router

import (
	"context"
	"fmt"
	"time"
)

// IdempotencyStore interface for idempotency tracking
type IdempotencyStore interface {
	CheckAndSet(ctx context.Context, key string, ttl time.Duration) (bool, error)
}

// TopicManager manages Kafka topics
type TopicManager struct {
	producer *KafkaProducer
}

// NewTopicManager creates a new topic manager
func NewTopicManager(producer *KafkaProducer) *TopicManager {
	return &TopicManager{producer: producer}
}

// EnsureTopic ensures a topic exists (Kafka auto-creates topics in MVP)
func (tm *TopicManager) EnsureTopic(ctx context.Context, topic string) error {
	// In MVP, Kafka auto-creates topics
	// In production, use AdminClient to create topics explicitly
	return nil
}

// GetTopicName returns the full topic name with tenant prefix
func GetTopicName(tenantID, topic string) string {
	return fmt.Sprintf("acb.%s.%s", tenantID, topic)
}


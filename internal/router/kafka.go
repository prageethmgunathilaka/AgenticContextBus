//go:build cgo
// +build cgo

package router

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/acb/internal/models"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

// KafkaProducer wraps Kafka producer
type KafkaProducer struct {
	producer *kafka.Producer
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(bootstrapServers string) (*KafkaProducer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"client.id":         "acb-producer",
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &KafkaProducer{producer: producer}, nil
}

// Produce sends a message to Kafka topic
func (p *KafkaProducer) Produce(ctx context.Context, topic string, message *models.Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Determine partition key
	partitionKey := message.To
	if partitionKey == "" {
		partitionKey = message.From // Use sender for broadcast
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(partitionKey),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "message-id", Value: []byte(message.ID)},
			{Key: "from", Value: []byte(message.From)},
			{Key: "type", Value: []byte(string(message.Type))},
		},
	}

	deliveryChan := make(chan kafka.Event)
	err = p.producer.Produce(msg, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	// Wait for delivery with timeout
	select {
	case e := <-deliveryChan:
		msg := e.(*kafka.Message)
		if msg.TopicPartition.Error != nil {
			return fmt.Errorf("delivery failed: %w", msg.TopicPartition.Error)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return fmt.Errorf("delivery timeout")
	}
}

// Close closes the producer
func (p *KafkaProducer) Close() {
	p.producer.Close()
}

// KafkaConsumer wraps Kafka consumer
type KafkaConsumer struct {
	consumer *kafka.Consumer
	topics   []string
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(bootstrapServers, groupID string, topics []string) (*KafkaConsumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	return &KafkaConsumer{
		consumer: consumer,
		topics:   topics,
	}, nil
}

// Consume consumes messages from Kafka
func (c *KafkaConsumer) Consume(ctx context.Context, handler func(*models.Message) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				return fmt.Errorf("failed to read message: %w", err)
			}

			var message models.Message
			if err := json.Unmarshal(msg.Value, &message); err != nil {
				// Log error but continue
				continue
			}

			if err := handler(&message); err != nil {
				// Log error but continue
				continue
			}

            // Commit offset
            if _, err := c.consumer.CommitMessage(msg); err != nil {
                return fmt.Errorf("failed to commit message: %w", err)
            }
		}
	}
}

// Close closes the consumer
func (c *KafkaConsumer) Close() error {
	return c.consumer.Close()
}

// Router handles message routing
type Router struct {
	producer *KafkaProducer
}

// NewRouter creates a new message router
func NewRouter(producer *KafkaProducer) *Router {
	return &Router{
		producer: producer,
	}
}

// SendTo sends a message to a specific agent
func (r *Router) SendTo(ctx context.Context, toAgentID string, topic string, message *models.Message) error {
	if message.ID == "" {
		message.ID = uuid.New().String()
	}
	if message.IdempotencyKey == "" {
		message.IdempotencyKey = uuid.New().String()
	}
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}

	message.To = toAgentID
	message.Topic = topic

	// Determine topic name (tenant-specific)
	kafkaTopic := fmt.Sprintf("acb.default.%s", topic) // MVP: single tenant

	return r.producer.Produce(ctx, kafkaTopic, message)
}

// Broadcast sends a message to all agents
func (r *Router) Broadcast(ctx context.Context, topic string, message *models.Message) error {
	message.To = "" // Empty for broadcast
	return r.SendTo(ctx, "", topic, message)
}

// Request sends a request message and waits for reply
func (r *Router) Request(ctx context.Context, toAgentID string, topic string, message *models.Message, timeout time.Duration) (*models.Message, error) {
	// Generate correlation ID
	message.CorrelationID = uuid.New().String()
	message.ReplyTo = fmt.Sprintf("reply.%s", uuid.New().String())

	// Send request
	if err := r.SendTo(ctx, toAgentID, topic, message); err != nil {
		return nil, err
	}

	// TODO: Wait for reply (implement reply queue)
	// For MVP, return immediately
	return nil, fmt.Errorf("request-reply not fully implemented in MVP")
}

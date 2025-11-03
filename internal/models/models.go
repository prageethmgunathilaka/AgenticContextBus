package models

import (
	"time"
)

// AgentStatus represents the current status of an agent
type AgentStatus string

const (
	AgentStatusOnline  AgentStatus = "online"
	AgentStatusOffline AgentStatus = "offline"
	AgentStatusUnknown AgentStatus = "unknown"
)

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeEvent    MessageType = "event"    // Fire-and-forget event
	MessageTypeCommand  MessageType = "command"  // Request for action
	MessageTypeQuery    MessageType = "query"    // Request for data
	MessageTypeResponse MessageType = "response" // Reply to query/command
)

// ContextScope defines who can access a context
type ContextScope string

const (
	ScopePublic  ContextScope = "public"  // All agents in tenant can access
	ScopePrivate ContextScope = "private" // Only creating agent can access
	ScopeGroup   ContextScope = "group"   // Specific group of agents
	ScopeShared  ContextScope = "shared"  // Explicitly shared with specific agents
)

// Agent represents an autonomous agent in the system
type Agent struct {
	ID           string            `json:"id" db:"id"`                     // Unique agent identifier
	Type         string            `json:"type" db:"type"`                 // Agent type (e.g., "ml", "rpa", "chatbot")
	Location     string            `json:"location" db:"location"`         // Physical/logical location
	Capabilities []string          `json:"capabilities" db:"capabilities"` // What this agent can do
	Metadata     map[string]string `json:"metadata" db:"metadata"`         // Custom metadata
	Status       AgentStatus       `json:"status" db:"status"`             // Current status
	TenantID     string            `json:"tenant_id" db:"tenant_id"`       // Tenant ownership
	CreatedAt    time.Time         `json:"created_at" db:"created_at"`
	LastSeen     time.Time         `json:"last_seen" db:"last_seen"` // Last heartbeat
}

// StorageRef points to large payloads in object storage (Phase 2+)
type StorageRef struct {
	Backend string `json:"backend"` // "s3", "minio", "azure", etc.
	Bucket  string `json:"bucket"`  // Bucket/container name
	Key     string `json:"key"`     // Object key/path
	Size    int64  `json:"size"`    // Size in bytes
}

// AccessControl defines who can access a context
type AccessControl struct {
	Scope      ContextScope `json:"scope"`                 // Access scope
	AllowedIDs []string     `json:"allowed_ids,omitempty"` // Specific agent IDs (for group/shared scopes)
}

// Context represents shared state/knowledge between agents
type Context struct {
	ID            string            `json:"id"`                                 // Unique context identifier
	Type          string            `json:"type"`                               // Context type (e.g., "user-profile", "model-weights")
	AgentID       string            `json:"agent_id" db:"agent_id"`             // Agent that created this context
	TenantID      string            `json:"tenant_id" db:"tenant_id"`           // Tenant ownership
	Payload       []byte            `json:"payload,omitempty"`                  // Actual data (for small contexts <1MB)
	PayloadRef    *StorageRef       `json:"payload_ref,omitempty"`              // Reference to large payload in object storage
	Metadata      map[string]string `json:"metadata" db:"metadata"`             // Additional metadata
	Version       string            `json:"version" db:"version"`               // Schema version
	SchemaID      string            `json:"schema_id,omitempty"`                // Schema registry ID
	TTL           time.Duration     `json:"ttl" db:"ttl"`                       // Time to live
	AccessControl AccessControl     `json:"access_control" db:"access_control"` // Who can access
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
	ExpiresAt     time.Time         `json:"expires_at" db:"expires_at"`       // Auto-delete after this
	Checksum      string            `json:"checksum,omitempty" db:"checksum"` // SHA-256 for integrity
}

// Message represents communication between agents
type Message struct {
	ID             string            `json:"id"`                       // Unique message ID
	From           string            `json:"from"`                     // Sender agent ID
	To             string            `json:"to,omitempty"`             // Recipient agent ID (empty for broadcast)
	Topic          string            `json:"topic"`                    // Kafka topic / routing key
	ContextID      string            `json:"context_id,omitempty"`     // Referenced context
	Context        *Context          `json:"context,omitempty"`        // Embedded context (for small data)
	Type           MessageType       `json:"type"`                     // Message type
	IdempotencyKey string            `json:"idempotency_key"`          // For deduplication
	CorrelationID  string            `json:"correlation_id,omitempty"` // For request-reply
	ReplyTo        string            `json:"reply_to,omitempty"`       // Reply address
	Metadata       map[string]string `json:"metadata"`                 // Custom metadata
	Timestamp      time.Time         `json:"timestamp"`                // When sent
	ExpiresAt      *time.Time        `json:"expires_at,omitempty"`     // Optional expiration
}

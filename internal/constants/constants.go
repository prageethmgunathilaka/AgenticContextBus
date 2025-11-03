package constants

// Default configuration constants
const (
	// Streaming
	DefaultChunkSize = 1 * 1024 * 1024 // 1MB
	MaxChunkSize     = 1 * 1024 * 1024 // 1MB max

	// Context size limits
	MaxDirectContextSize    = 1 * 1024 * 1024   // 1MB
	MaxStreamingContextSize = 100 * 1024 * 1024 // 100MB

	// Rate limits (per agent)
	DefaultRateLimitRequests    = 1000              // requests per minute
	DefaultRateLimitBandwidth   = 100 * 1024 * 1024 // 100 MB per minute
	DefaultMaxConcurrentStreams = 10

	// Token expiration
	DefaultAccessTokenTTL  = 3600   // 1 hour in seconds
	DefaultRefreshTokenTTL = 604800 // 7 days in seconds

	// Agent heartbeat
	DefaultHeartbeatInterval = 30 // seconds
	AgentTimeoutThreshold    = 90 // seconds (mark offline after 90s without heartbeat)

	// Storage
	DefaultContextTTL = 86400 // 24 hours in seconds

	// Idempotency
	IdempotencyKeyTTL = 86400 // 24 hours in seconds

	// Database
	DefaultMaxConnections = 25
	DefaultMinConnections = 5

	// Kafka
	DefaultKafkaReplicationFactor = 1
	DefaultKafkaPartitions        = 3

	// Timeouts
	DefaultHTTPTimeout = 30 // seconds
	DefaultDialTimeout = 10 // seconds

	// Pagination
	DefaultPageLimit  = 100
	MaxPageLimit      = 1000
	DefaultPageOffset = 0
)

// Tenant resource limits (for Phase 2+)
const (
	FreeTierMaxAgents   = 3
	FreeTierMaxMessages = 10000
	FreeTierMaxStorage  = 1 * 1024 * 1024 * 1024 // 1GB

	StarterTierMaxAgents   = 10
	StarterTierMaxMessages = 500000
	StarterTierMaxStorage  = 10 * 1024 * 1024 * 1024 // 10GB

	ProfessionalTierMaxAgents   = 50
	ProfessionalTierMaxMessages = 5000000
	ProfessionalTierMaxStorage  = 100 * 1024 * 1024 * 1024 // 100GB
)

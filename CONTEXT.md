# ACB Project Context for AI Assistants

## Project Overview
**Agentic Context Bus (ACB)** is an enterprise-grade communication framework for distributed AI agents to share context across remote locations. This document provides complete context for any AI assistant working on this project.

## Key Requirements
- **Scale**: 10-50 agents initially, designed to scale to 1000s
- **Performance**: P99 latency <50ms (small messages), 100+ MB/s streaming throughput
- **Security**: TLS 1.3, JWT authentication, RBAC, comprehensive audit logging
- **Deployment**: Production-ready from day one with Docker/Kubernetes support
- **Developer**: Learning Go while building (provide helpful comments and explanations)

## Core Problem Being Solved
Allow autonomous agents (AI, RPA, automation, microservices) running in different locations to:
- Share context (state, memory, knowledge, events)
- Handle any size data: small messages to GB-scale streams
- Discover and communicate with each other
- Maintain security and observability
- Scale from prototype to enterprise production

## Architecture Decisions

### 1. Hub-Based Architecture (Not Peer-to-Peer)
**Decision**: All agent communication goes through central ACB hub

**Rationale**:
- Simplicity: Faster to build and reason about
- Complete observability: See all messages for debugging/monitoring
- Centralized security: Single point for auth, audit, policy enforcement
- Easier multi-tenancy: Logical isolation per customer

**Future**: A2A (Agent-to-Agent) direct communication can be added in Phase 2+ without breaking changes. SDK designed to support it transparently.

**Consequences**:
- Hub is potential bottleneck (mitigate with horizontal scaling)
- Slightly higher latency vs direct P2P (acceptable for MVP)

### 2. Message Delivery: At-Least-Once
**Decision**: Messages guaranteed to arrive, but may duplicate

**Rationale**:
- Right balance between reliability and complexity
- Kafka naturally supports this
- Simpler than exactly-once semantics
- Sufficient for 95% of use cases

**Implementation**:
- Agents must handle idempotency (use idempotency keys)
- SDK provides idempotency key generation
- Server tracks recent keys to detect duplicates

**Alternatives Considered**:
- At-most-once: Too risky (message loss)
- Exactly-once: Too complex for MVP, can add later

### 3. Hybrid Consistency Model
**Decision**: Mix of strong and eventual consistency based on data type

**Implementation**:
- **Strong consistency** (PostgreSQL with transactions):
  - Agent registration/state
  - Permissions and access control
  - Billing/usage data
  - Critical configuration
  
- **Eventual consistency** (Kafka streaming):
  - Messages between agents
  - Events and logs
  - Metrics/telemetry
  - Non-critical context updates

**Rationale**: Use the right consistency model for each use case. Strong consistency where correctness matters, eventual where speed matters.

### 4. Context Handling Strategy
**Decision**: Three-tier approach based on context size

**Implementation**:
- **<1MB**: Direct via Kafka messages
  - Fast, low overhead
  - Serialized in message payload
  
- **1-100MB**: Streaming with chunking (Strategy A)
  - gRPC bidirectional streaming
  - 1MB chunks
  - Progress tracking
  - Resume capability
  - Checksum validation (SHA-256)
  
- **>100MB**: Document for future (Phase 2+)
  - Object storage (S3/MinIO) with reference in message
  - Lazy loading when needed
  - CDN for global distribution

**Rationale**: Optimize for different data sizes. Don't force large files through message bus, but don't require object storage for medium files.

### 5. Go as Primary Language
**Decision**: Build ACB core in Go

**Rationale**:
- **Performance**: Compiled, fast, low latency, predictable GC
- **Concurrency**: Goroutines perfect for handling 1000s of agents
- **Streaming**: Native io.Reader/Writer, excellent for our use case
- **Deployment**: Single binary, no runtime dependencies
- **Cloud-native**: First-class support in Kubernetes/Docker ecosystem
- **Kafka integration**: Excellent client libraries (Confluent, Shopify Sarama)
- **Learning opportunity**: Developer wants to learn Go

**Developer Context**: This is a learning project. Code should:
- Include helpful comments explaining Go concepts
- Compare to Java/Kotlin equivalents where relevant
- Use idiomatic Go patterns
- Progressively introduce concepts (start simple)

**Future**: Multi-language SDKs (Python, Java, JavaScript) will be created, but core platform is Go.

### 6. Security Model
**Authentication**: 
- JWT tokens (primary method)
- API keys for automation
- mTLS support for enterprise (optional, Phase 2+)

**Authorization**: 
- RBAC with predefined roles:
  - `admin`: Full access, manage all tenants
  - `agent-producer`: Can send messages and create contexts
  - `agent-consumer`: Can receive messages and read contexts
  - `agent-full`: Both producer and consumer
  - `observer`: Read-only access for monitoring

**Encryption**:
- TLS 1.3 mandatory for all connections (no plaintext!)
- Optional E2E encryption for sensitive contexts (agent encrypts, ACB can't see)

**Audit**: 
- Log ALL security events:
  - Authentication attempts (success/failure)
  - Authorization decisions
  - Context access (read/write)
  - Permission changes
  - Agent registration/deregistration
- Store in PostgreSQL with 90-day retention minimum
- Exportable for compliance (GDPR, SOC2)

### 7. Multi-Tenancy Approach
**Decision**: Logical isolation with optional physical isolation for enterprise

**Isolation Strategy**:
- **Kafka**: Separate topics per tenant
  - Pattern: `acb.{tenant-id}.{context-type}`
  - Example: `acb.company-abc.agent-events`
  
- **Redis**: Namespace prefixes
  - Pattern: `tenant:{tenant-id}:{key}`
  - Example: `tenant:company-abc:agent:agent-123`
  
- **PostgreSQL**: Row-level security
  - Every table has `tenant_id` column
  - All queries filtered by tenant
  - Foreign key constraints enforce isolation
  
- **Network**: Optional VPC isolation for enterprise tier

**Provisioning**: Auto-provision new tenant on signup (<60 seconds)

**Quotas**: Per-tenant resource limits enforced:
- Max agents
- Messages per minute
- Bandwidth (MB/minute)
- Storage (GB)

## Technology Stack

### Core
- **Language**: Go 1.21+
- **Message Bus**: Apache Kafka 3.x (Confluent client library: `github.com/confluentinc/confluent-kafka-go`)
- **Hot Storage**: Redis 7+ with TLS (`github.com/redis/go-redis`)
- **Warm Storage**: PostgreSQL 15+ (`github.com/jackc/pgx`)
- **Streaming**: gRPC (`google.golang.org/grpc`) + WebSocket (gorilla/websocket)
- **HTTP Router**: Gin or Chi (to be decided)

### Security
- **JWT**: `github.com/golang-jwt/jwt`
- **TLS**: Go standard library crypto/tls
- **Hashing**: bcrypt for passwords, SHA-256 for checksums

### Observability
- **Metrics**: Prometheus (`github.com/prometheus/client_golang`)
- **Tracing**: OpenTelemetry → Jaeger
- **Logging**: zerolog or zap (structured JSON logs)
- **Visualization**: Grafana

### Development
- **Containerization**: Docker, Docker Compose
- **Orchestration**: Kubernetes
- **Testing**: Go testing package, testify
- **Mocking**: gomock or mockery

## Project Structure

```
acb/
├── cmd/                           # Executable entry points
│   ├── acb-server/               # Main ACB server
│   │   └── main.go
│   ├── acb-cli/                  # CLI tool
│   │   └── main.go
│   └── acb-agent-demo/           # Demo agents
│       ├── hello-world/
│       └── streaming-demo/
├── internal/                      # Private application code (not importable)
│   ├── auth/                     # Authentication/authorization
│   │   ├── jwt.go
│   │   ├── rbac.go
│   │   └── middleware.go
│   ├── registry/                 # Agent registry
│   │   ├── service.go
│   │   ├── store.go
│   │   └── models.go
│   ├── context/                  # Context management
│   │   ├── manager.go
│   │   ├── lifecycle.go
│   │   └── versioning.go
│   ├── router/                   # Message routing
│   │   ├── kafka.go
│   │   ├── router.go
│   │   └── dlq.go
│   ├── stream/                   # Streaming handlers
│   │   ├── grpc.go
│   │   ├── websocket.go
│   │   └── chunker.go
│   ├── storage/                  # Database abstractions
│   │   ├── postgres.go
│   │   ├── redis.go
│   │   └── interfaces.go
│   ├── tenant/                   # Multi-tenancy
│   │   ├── manager.go
│   │   └── isolation.go
│   ├── billing/                  # Usage metering
│   │   ├── meter.go
│   │   └── aggregator.go
│   ├── metrics/                  # Prometheus metrics
│   │   └── metrics.go
│   ├── health/                   # Health checks
│   │   └── checker.go
│   └── server/                   # HTTP/gRPC servers
│       ├── grpc.go
│       ├── http.go
│       └── tls.go
├── pkg/                          # Public libraries (importable)
│   └── acb-sdk/                  # Agent SDK
│       ├── client.go
│       ├── types.go
│       ├── stream.go
│       ├── builder.go
│       └── errors.go
├── api/                          # API definitions
│   ├── proto/                    # gRPC protobuf definitions
│   │   ├── registry.proto
│   │   ├── context.proto
│   │   └── stream.proto
│   └── openapi/                  # REST OpenAPI specs
│       └── acb-api.yaml
├── deployments/                   # Deployment configs
│   ├── docker/                   # Dockerfiles
│   │   ├── Dockerfile.server
│   │   ├── Dockerfile.cli
│   │   └── Dockerfile.agent-demo
│   └── k8s/                      # Kubernetes manifests
│       ├── namespace.yaml
│       ├── deployment.yaml
│       ├── service.yaml
│       └── ingress.yaml
├── docs/                         # Documentation
│   ├── ARCHITECTURE.md           # Detailed architecture
│   ├── ADR/                      # Architecture Decision Records
│   │   ├── 001-hub-vs-p2p.md
│   │   ├── 002-kafka-over-rabbitmq.md
│   │   ├── 003-at-least-once-delivery.md
│   │   ├── 004-go-as-primary-language.md
│   │   └── 005-multi-tenancy-approach.md
│   ├── quickstart.md             # 5-minute tutorial
│   ├── guides/                   # User guides
│   │   ├── authentication.md
│   │   ├── streaming.md
│   │   ├── deployment.md
│   │   └── monitoring.md
│   ├── sdk-reference/            # SDK API docs
│   │   └── go.md
│   └── api-reference/            # gRPC/REST API docs
│       ├── grpc.md
│       └── rest.md
├── examples/                      # Example agents
│   ├── hello-world/
│   │   ├── agent-a/
│   │   └── agent-b/
│   ├── streaming-demo/
│   └── multi-agent/
├── tests/                        # Integration tests
│   ├── integration/
│   └── e2e/
├── scripts/                      # Utility scripts
│   ├── setup-dev.sh
│   ├── generate-certs.sh
│   └── run-tests.sh
├── .github/                      # GitHub configs
│   └── workflows/
│       ├── ci.yml
│       └── release.yml
├── CONTEXT.md                    # This file - complete context
├── README.md                     # Project overview
├── DEVELOPMENT.md                # Development setup
├── LICENSE                       # Apache 2.0 or Business Source License
├── .gitignore
├── go.mod                        # Go module definition
├── go.sum                        # Go dependencies
├── Makefile                      # Build automation
├── docker-compose.yml            # Local development environment
└── docker-compose.dev.yml        # Additional dev services
```

## Core Data Models

### Agent
```go
// Agent represents an autonomous agent in the system
type Agent struct {
    ID           string                 `json:"id" db:"id"`                     // Unique agent identifier
    Type         string                 `json:"type" db:"type"`                 // Agent type (e.g., "ml", "rpa", "chatbot")
    Location     string                 `json:"location" db:"location"`         // Physical/logical location
    Capabilities []string               `json:"capabilities" db:"capabilities"` // What this agent can do
    Metadata     map[string]string      `json:"metadata" db:"metadata"`        // Custom metadata
    Status       AgentStatus            `json:"status" db:"status"`            // Current status
    TenantID     string                 `json:"tenant_id" db:"tenant_id"`      // Tenant ownership
    CreatedAt    time.Time              `json:"created_at" db:"created_at"`
    LastSeen     time.Time              `json:"last_seen" db:"last_seen"`      // Last heartbeat
}

type AgentStatus string

const (
    AgentStatusOnline  AgentStatus = "online"
    AgentStatusOffline AgentStatus = "offline"
    AgentStatusUnknown AgentStatus = "unknown"
)
```

### Context
```go
// Context represents shared state/knowledge between agents
type Context struct {
    ID            string            `json:"id"`                          // Unique context identifier
    Type          string            `json:"type"`                        // Context type (e.g., "user-profile", "model-weights")
    AgentID       string            `json:"agent_id"`                    // Agent that created this context
    TenantID      string            `json:"tenant_id"`                   // Tenant ownership
    Payload       []byte            `json:"payload,omitempty"`           // Actual data (for small contexts)
    PayloadRef    *StorageRef       `json:"payload_ref,omitempty"`       // Reference to large payload in object storage
    Metadata      map[string]string `json:"metadata"`                    // Additional metadata
    Version       string            `json:"version"`                     // Schema version
    SchemaID      string            `json:"schema_id,omitempty"`         // Schema registry ID
    TTL           time.Duration     `json:"ttl"`                         // Time to live
    AccessControl AccessControl     `json:"access_control"`              // Who can access
    CreatedAt     time.Time         `json:"created_at"`
    ExpiresAt     time.Time         `json:"expires_at"`                  // Auto-delete after this
    Checksum      string            `json:"checksum,omitempty"`          // SHA-256 for integrity
}

// StorageRef points to large payloads in object storage
type StorageRef struct {
    Backend  string `json:"backend"`  // "s3", "minio", "azure", etc.
    Bucket   string `json:"bucket"`   // Bucket/container name
    Key      string `json:"key"`      // Object key/path
    Size     int64  `json:"size"`     // Size in bytes
}

// AccessControl defines who can access this context
type AccessControl struct {
    Scope      ContextScope `json:"scope"`                  // Access scope
    AllowedIDs []string     `json:"allowed_ids,omitempty"`  // Specific agent IDs (for group/shared)
}

type ContextScope string

const (
    ScopePublic  ContextScope = "public"  // All agents in tenant can access
    ScopePrivate ContextScope = "private" // Only creating agent can access
    ScopeGroup   ContextScope = "group"   // Specific group of agents
    ScopeShared  ContextScope = "shared"  // Explicitly shared with specific agents
)
```

### Message
```go
// Message represents communication between agents
type Message struct {
    ID             string            `json:"id"`                        // Unique message ID
    From           string            `json:"from"`                      // Sender agent ID
    To             string            `json:"to,omitempty"`              // Recipient agent ID (empty for broadcast)
    Topic          string            `json:"topic"`                     // Kafka topic / routing key
    ContextID      string            `json:"context_id,omitempty"`      // Referenced context
    Context        *Context          `json:"context,omitempty"`         // Embedded context (for small data)
    Type           MessageType       `json:"type"`                      // Message type
    IdempotencyKey string            `json:"idempotency_key"`           // For deduplication
    CorrelationID  string            `json:"correlation_id,omitempty"`  // For request-reply
    ReplyTo        string            `json:"reply_to,omitempty"`        // Reply address
    Metadata       map[string]string `json:"metadata"`                  // Custom metadata
    Timestamp      time.Time         `json:"timestamp"`                 // When sent
    ExpiresAt      *time.Time        `json:"expires_at,omitempty"`      // Optional expiration
}

type MessageType string

const (
    MessageTypeEvent    MessageType = "event"    // Fire-and-forget event
    MessageTypeCommand  MessageType = "command"  // Request for action
    MessageTypeQuery    MessageType = "query"    // Request for data
    MessageTypeResponse MessageType = "response" // Reply to query/command
)
```

## API Design

### Agent SDK (Go) - Key Methods

```go
// Client creation with builder pattern
client := acb.NewClient(
    acb.WithEndpoint("acb.example.com:443"),
    acb.WithCredentials(jwt),
    acb.WithTLS(tlsConfig),
    acb.WithTimeout(30*time.Second),
)
defer client.Close()

// Simple context sharing
err := client.ShareContext(ctx, "user-profile", userData)

// With options
err := client.ShareContext(ctx, "execution-state", state,
    acb.WithScope(acb.ScopeGroup),
    acb.WithTTL(24*time.Hour),
    acb.WithMetadata(map[string]string{"priority": "high", "version": "1.0"}),
)

// Subscribe to context updates
subscription := client.Subscribe("agent-events", func(ctx context.Context, c *acb.Context) error {
    log.Printf("Received context: %s", c.ID)
    return handleContext(c)
})
defer subscription.Unsubscribe()

// Subscribe with filtering
subscription := client.Subscribe("agent-events",
    func(ctx context.Context, c *acb.Context) error {
        return handleContext(c)
    },
    acb.WithFilter("metadata.priority", "high"),
)

// Streaming large contexts (send)
err := client.StreamContext(ctx, "model-weights").
    FromReader(file).
    WithChunkSize(1024 * 1024). // 1MB chunks
    OnProgress(func(progress float64) {
        log.Printf("Upload progress: %.1f%%", progress*100)
    }).
    Send()

// Streaming large contexts (receive)
reader, err := client.ReceiveContextStream(ctx, contextID)
if err != nil {
    log.Fatal(err)
}
defer reader.Close()

written, err := io.Copy(outputFile, reader)

// Request-reply pattern (synchronous)
response, err := client.Request(ctx, "agent-123", "get-status").Wait()

// Request-reply pattern (asynchronous)
future := client.Request(ctx, "agent-123", "get-status")
// Do other work...
response, err := future.Wait()

// Broadcast to all agents
err := client.Broadcast(ctx, "system-announcement", announcement)

// Send to specific agent
err := client.SendTo(ctx, "agent-456", "task-assignment", task)
```

### gRPC Service Definitions

```protobuf
// Registry service for agent management
service AgentRegistry {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Unregister(UnregisterRequest) returns (UnregisterResponse);
  rpc Heartbeat(HeartbeatRequest) returns (HeartbeatResponse);
  rpc Discover(DiscoverRequest) returns (DiscoverResponse);
  rpc GetAgent(GetAgentRequest) returns (GetAgentResponse);
  rpc ListAgents(ListAgentsRequest) returns (ListAgentsResponse);
}

// Context service for context management
service ContextService {
  rpc CreateContext(CreateContextRequest) returns (CreateContextResponse);
  rpc GetContext(GetContextRequest) returns (GetContextResponse);
  rpc UpdateContext(UpdateContextRequest) returns (UpdateContextResponse);
  rpc DeleteContext(DeleteContextRequest) returns (DeleteContextResponse);
  rpc ListContexts(ListContextsRequest) returns (ListContextsResponse);
  rpc Subscribe(SubscribeRequest) returns (stream ContextEvent);
}

// Stream service for large context transfers
service StreamService {
  rpc StreamContext(stream ContextChunk) returns (StreamResponse);
  rpc ReceiveContext(ContextRequest) returns (stream ContextChunk);
  rpc GetStreamProgress(StreamProgressRequest) returns (StreamProgressResponse);
}
```

## Development Workflow

### Setting Up Local Environment

```bash
# Clone repository
git clone https://github.com/[username]/acb.git
cd acb

# Start all dependencies (Kafka, Redis, PostgreSQL)
make dev-up

# Run ACB server
make run-server

# In another terminal, run demo agent
make run-demo

# Run tests
make test

# Stop all services
make dev-down
```

### Go Learning Path for This Project

**Week 1: Basics**
- Syntax, types, structs, functions
- Error handling
- Packages and imports
- Will implement: Basic types, project structure

**Week 2: Concurrency**
- Goroutines
- Channels
- Context package
- sync package (WaitGroup, Mutex)
- Will implement: Message router, concurrent handlers

**Week 3: Networking & I/O**
- HTTP servers
- gRPC
- io.Reader/Writer
- JSON encoding/decoding
- Will implement: Streaming, SDK

**Week 4+: Production**
- Testing
- Profiling
- Error wrapping
- Dependency injection
- Will implement: Security, monitoring, polish

### Code Style Conventions

**Naming**:
- Packages: lowercase, single word (`auth`, `registry`)
- Files: lowercase with underscores (`agent_registry.go`)
- Types: PascalCase (`AgentRegistry`)
- Functions: PascalCase for exported, camelCase for internal
- Interfaces: Single method = verb+er (`Reader`, `Writer`), multi-method = noun

**Error Handling**:
```go
// Always check errors
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doing something: %w", err) // Wrap errors with context
}

// Use errors.Is for comparison
if errors.Is(err, ErrNotFound) {
    // Handle not found
}

// Use errors.As for type assertion
var validationErr *ValidationError
if errors.As(err, &validationErr) {
    // Handle validation error
}
```

**Concurrency**:
```go
// Use goroutines for concurrent work
go processMessage(msg)

// Use channels for communication
messages := make(chan Message, 100)

// Always use context for cancellation
ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
defer cancel()

// Clean up goroutines
done := make(chan struct{})
go func() {
    defer close(done)
    // Do work
}()
<-done // Wait for completion
```

## Performance Requirements

### Latency Targets
- **P50**: <10ms for small messages (<1KB)
- **P95**: <30ms
- **P99**: <50ms
- **Streaming throughput**: 100+ MB/s aggregate

### Scalability Targets
- **Phase 1 (MVP)**: 50 concurrent agents
- **Phase 2**: 500 concurrent agents
- **Phase 3**: 1000+ concurrent agents
- Design for horizontal scaling (stateless services)

### Resource Limits (Per Agent)
- **Messages**: 1000/minute (configurable per tier)
- **Bandwidth**: 100 MB/minute
- **Concurrent streams**: 10
- **Storage**: 1GB (can be adjusted)

## Security Requirements

### Authentication
- JWT tokens with 1-hour expiration
- Refresh token mechanism (7-day expiration)
- API keys for automation/service accounts
- mTLS optional for enterprise

### Authorization
- RBAC enforced at every API endpoint
- Context-level access control
- Tenant isolation strictly enforced

### Audit Requirements
- Log all authentication attempts
- Log all context read/write operations
- Log all permission changes
- Retention: 90 days minimum
- Exportable in JSON format

## Monitoring & Observability

### Key Metrics (Prometheus)
```
# System metrics
acb_messages_total{type, tenant_id}
acb_message_latency_seconds{operation}
acb_active_agents{tenant_id}
acb_context_storage_bytes{tenant_id}
acb_errors_total{type, operation}

# Agent metrics
acb_agent_messages_sent_total{agent_id}
acb_agent_messages_received_total{agent_id}
acb_agent_bandwidth_bytes{agent_id, direction}

# Performance metrics
acb_kafka_lag{topic}
acb_redis_operations_total{operation}
acb_db_query_duration_seconds{query}
```

### Logging Standards
```go
// Use structured logging with zerolog or zap
log.Info().
    Str("agent_id", agentID).
    Str("tenant_id", tenantID).
    Str("type", agentType).
    Msg("agent registered")

// Include correlation IDs
log.Error().
    Err(err).
    Str("correlation_id", correlationID).
    Str("from", fromAgent).
    Str("to", toAgent).
    Msg("failed to send message")
```

## SaaS Business Model

### Pricing Tiers
1. **Free**: 3 agents, 10k msgs/month, 1GB storage, community support
2. **Starter** ($99/mo): 10 agents, 500k msgs/month, 10GB storage, email support
3. **Professional** ($499/mo): 50 agents, 5M msgs/month, 100GB storage, priority support
4. **Enterprise** (custom): Unlimited agents, unlimited messages, dedicated infrastructure

### Usage Tracking
- Real-time metering: messages, bandwidth, storage
- Hourly aggregation for billing
- Usage dashboards for customers
- Alerts when approaching limits

## Testing Strategy

### Unit Tests
```go
func TestAgentRegistry_Register(t *testing.T) {
    tests := []struct {
        name    string
        agent   *Agent
        wantErr bool
    }{
        {
            name: "valid agent",
            agent: &Agent{ID: "agent-1", Type: "ml"},
            wantErr: false,
        },
        {
            name: "duplicate agent",
            agent: &Agent{ID: "agent-1", Type: "ml"},
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests
- Test with real Kafka, Redis, PostgreSQL (Docker Compose)
- Test complete flows

### End-to-End Tests
- Test full agent lifecycle
- Test streaming large contexts
- Test failure scenarios

## Open Questions / Future Considerations

1. **Object Storage**: When to implement for >100MB contexts? (Phase 2+)
2. **Schema Registry**: Use Confluent Schema Registry or build custom?
3. **Multi-Language SDKs**: Priority order (Python first, then Java, then JavaScript?)
4. **A2A Direct**: When to implement direct agent-to-agent communication? (Phase 2+)
5. **Agent Marketplace**: Build platform for pre-built agents?
6. **Workflow Orchestration**: Add DAG-based orchestration?

## Common Pitfalls to Avoid

### Go-Specific
- ❌ Not checking errors (every error must be handled!)
- ❌ Goroutine leaks (always clean up)
- ❌ Not using context for cancellation
- ❌ Race conditions (use channels or mutexes)
- ✅ "Don't communicate by sharing memory, share memory by communicating"

### Architecture
- ❌ Over-engineering early (YAGNI - You Ain't Gonna Need It)
- ❌ Premature optimization
- ❌ Not planning for failures
- ✅ Start simple, iterate based on real usage

## Success Criteria

### Phase 1 - MVP (Weeks 1-4)
- [ ] 2 agents can exchange messages
- [ ] Agent discovery functional
- [ ] JWT authentication working
- [ ] Stream 10MB+ contexts successfully
- [ ] Docker Compose dev environment
- [ ] Basic documentation

### Phase 2 - Production Ready (Weeks 5-6)
- [ ] 50 concurrent agents supported
- [ ] P99 latency <50ms
- [ ] TLS enforced
- [ ] Audit logging complete
- [ ] Monitoring dashboards
- [ ] Multi-tenancy working
- [ ] Comprehensive docs

### Phase 3 - SaaS Launch (Weeks 7-8)
- [ ] Marketing website live
- [ ] Self-service signup
- [ ] Usage metering
- [ ] Billing integration
- [ ] 10+ beta customers

## Additional Resources

- **Go Documentation**: https://go.dev/doc/
- **Kafka Documentation**: https://kafka.apache.org/documentation/
- **gRPC Go Tutorial**: https://grpc.io/docs/languages/go/
- **Go Best Practices**: https://go.dev/doc/effective_go

---

**Last Updated**: 2025-10-31  
**For Questions**: Refer to `/docs/` directory for detailed documentation  
**For Implementation**: Follow the todos in order, marking as in_progress when starting


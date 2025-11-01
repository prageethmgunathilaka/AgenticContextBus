# Product Requirements Document (PRD)
## Agentic Context Bus (ACB) - MVP / Phase 1

**Version**: 1.0  
**Date**: 2025-01-XX  
**Status**: Draft

---

## Executive Summary

The **Agentic Context Bus (ACB)** is an enterprise-grade communication framework that enables distributed AI agents to securely share context, state, and knowledge across remote locations. This PRD defines the Minimum Viable Product (MVP) for Phase 1, designed to support 2-50 agents with production-ready architecture that scales to 1000s without rework.

### Key Objectives
- Enable autonomous agents to discover and communicate with each other
- Support context sharing from small messages (<1KB) to large streams (1-100MB)
- Provide enterprise-grade security (TLS 1.3, JWT, RBAC)
- Deliver production-ready system with Docker/Kubernetes support
- Achieve P99 latency <50ms and 100+ MB/s streaming throughput

### Target Users
- AI/ML agents requiring distributed context sharing
- RPA (Robotic Process Automation) agents
- Autonomous microservices
- Multi-agent systems

---

## Problem Statement & Solution

### Problem
Modern autonomous agents operate in distributed environments but lack a standardized, secure, and scalable way to:
- Share context (state, memory, knowledge, events) across locations
- Discover other agents dynamically
- Handle variable data sizes efficiently (from KB to GB)
- Maintain security and observability at scale

### Solution
ACB provides a centralized hub-based communication framework that:
- Routes all agent communication through a secure, observable hub
- Supports multiple data transfer strategies (direct messages, streaming, future object storage)
- Implements comprehensive security (authentication, authorization, audit logging)
- Scales horizontally from prototype to enterprise production

---

## MVP Scope (Phase 1)

### In Scope - Core Features

#### 1. Agent Management
- **Agent Registration**: Agents can register with ACB hub
- **Agent Discovery**: Agents can discover other agents by type, location, or capabilities
- **Heartbeat Mechanism**: Agents send periodic heartbeats to maintain online status
- **Agent Status Tracking**: Real-time status (online/offline/unknown)

#### 2. Context Sharing
- **Small Contexts (<1MB)**: Direct sharing via Kafka messages
- **Medium Contexts (1-100MB)**: Streaming via gRPC with chunking
- **Context Types**: Support for typed contexts (user-profile, model-weights, execution-state, etc.)
- **Context Metadata**: Custom metadata and versioning support
- **Access Control**: Public, private, group, and shared scopes

#### 3. Message Routing
- **Point-to-Point**: Send messages to specific agents
- **Broadcast**: Send messages to all agents in tenant
- **Topic-Based**: Route messages by topic/category
- **Request-Reply**: Synchronous and asynchronous request-reply patterns
- **At-Least-Once Delivery**: Messages guaranteed to arrive (may duplicate)

#### 4. Authentication & Authorization
- **JWT Authentication**: Primary authentication method with 1-hour expiration
- **Refresh Tokens**: 7-day refresh token mechanism
- **API Keys**: Support for automation/service accounts
- **RBAC**: Role-based access control (admin, agent-producer, agent-consumer, agent-full, observer)
- **Context-Level Access Control**: Fine-grained permissions per context

#### 5. Streaming (Large Contexts)
- **gRPC Bidirectional Streaming**: For contexts 1-100MB
- **Chunking**: 1MB chunks with progress tracking
- **Resume Capability**: Resume interrupted transfers
- **Checksum Validation**: SHA-256 for integrity verification
- **Progress Tracking**: Real-time upload/download progress

#### 6. Developer Experience
- **Go SDK**: Complete SDK with builder pattern
- **Docker Compose**: Local development environment
- **API Documentation**: Complete REST and gRPC API documentation
- **Demo Agents**: Example agents demonstrating key features

### Success Criteria - Phase 1

- ✅ 2 agents can exchange messages successfully
- ✅ Agent discovery functional (find agents by type/capabilities)
- ✅ JWT authentication working end-to-end
- ✅ Stream 10MB+ contexts successfully via gRPC
- ✅ Docker Compose dev environment operational
- ✅ Basic documentation complete (quickstart, API reference)

### Performance Targets - Phase 1

- **Latency**: P99 <50ms for small messages (<1KB)
- **Throughput**: 100+ MB/s aggregate streaming throughput
- **Scale**: Support 50 concurrent agents
- **Availability**: Single-node deployment (multi-node in Phase 2)

---

## Long-term Vision (Phase 2/3)

### Phase 2 - Production Ready (Weeks 5-6)
- Multi-tenancy with logical isolation
- Advanced security (mTLS, E2E encryption)
- Performance optimization for 500 concurrent agents
- Comprehensive monitoring and observability
- Production hardening (error handling, resilience)
- Load testing and optimization

### Phase 3 - SaaS Launch (Weeks 7-8)
- Self-service signup and tenant provisioning
- Usage metering and billing integration
- Marketing website and landing pages
- Customer dashboard and admin portal
- Beta customer onboarding
- Production launch preparation

**Note**: Architecture decisions in Phase 1 support Phase 2/3 without rework. Multi-tenancy is designed but not fully implemented in MVP (single tenant mode).

---

## Functional Requirements

### FR-1: Agent Registration
**Priority**: P0 (Critical)  
**Description**: Agents must be able to register with ACB hub providing:
- Agent ID (unique identifier)
- Agent type (ml, rpa, chatbot, etc.)
- Location (physical/logical location)
- Capabilities (what this agent can do)
- Metadata (custom key-value pairs)

**Acceptance Criteria**:
- Agent can register via REST API or gRPC
- Duplicate agent IDs are rejected
- Registration returns agent record with status
- Agent appears in discovery queries immediately after registration

### FR-2: Agent Discovery
**Priority**: P0 (Critical)  
**Description**: Agents must be able to discover other agents by:
- Agent type
- Location
- Capabilities
- All agents in tenant

**Acceptance Criteria**:
- Discovery queries return matching agents
- Results include agent status (online/offline)
- Results can be filtered and paginated
- Discovery works across all registered agents

### FR-3: Context Sharing (Small)
**Priority**: P0 (Critical)  
**Description**: Agents can share contexts <1MB directly via Kafka messages.

**Acceptance Criteria**:
- Context can be created with type, payload, metadata
- Context can be retrieved by ID
- Context can be updated
- Context can be deleted
- Access control enforced (public/private/group/shared)

### FR-4: Context Sharing (Large)
**Priority**: P0 (Critical)  
**Description**: Agents can stream contexts 1-100MB via gRPC with chunking.

**Acceptance Criteria**:
- Stream can be initiated with context metadata
- Data sent in 1MB chunks
- Progress tracked and reportable
- Stream can be resumed if interrupted
- Checksum validated on completion

### FR-5: Message Routing
**Priority**: P0 (Critical)  
**Description**: Agents can send messages:
- To specific agent (point-to-point)
- To all agents (broadcast)
- By topic (topic-based routing)

**Acceptance Criteria**:
- Messages delivered at-least-once
- Messages include idempotency keys for deduplication
- Request-reply pattern supported (sync and async)
- Messages can reference contexts

### FR-6: Authentication
**Priority**: P0 (Critical)  
**Description**: All API access requires authentication via:
- JWT tokens (primary)
- API keys (for automation)

**Acceptance Criteria**:
- JWT tokens validated on every request
- Tokens expire after 1 hour
- Refresh tokens allow renewal
- API keys authenticated
- Unauthenticated requests rejected

### FR-7: Authorization
**Priority**: P0 (Critical)  
**Description**: RBAC enforced at API level with roles:
- `admin`: Full access
- `agent-producer`: Can send messages and create contexts
- `agent-consumer`: Can receive messages and read contexts
- `agent-full`: Producer + consumer
- `observer`: Read-only access

**Acceptance Criteria**:
- Roles enforced at endpoints
- Context-level access control enforced
- Unauthorized requests rejected with appropriate error
- Audit log records all authorization decisions

### FR-8: Subscription (Context Updates)
**Priority**: P1 (High)  
**Description**: Agents can subscribe to context updates by topic.

**Acceptance Criteria**:
- Subscription established via gRPC streaming
- Updates delivered in real-time
- Filters supported (metadata-based)
- Subscription can be cancelled

---

## Non-Functional Requirements

### NFR-1: Performance
- **Latency**: P99 <50ms for small messages (<1KB)
- **Throughput**: 100+ MB/s aggregate streaming throughput
- **Concurrency**: Support 50 concurrent agents in Phase 1

### NFR-2: Security
- **TLS**: TLS 1.3 mandatory for all connections (no plaintext)
- **Encryption**: Data encrypted in transit
- **Audit**: All security events logged (authentication, authorization, access)
- **Retention**: Audit logs retained for 90 days minimum

### NFR-3: Reliability
- **Message Delivery**: At-least-once delivery guarantee
- **Idempotency**: SDK provides idempotency keys
- **Error Handling**: Graceful error handling with clear error messages
- **Health Checks**: Health check endpoints for monitoring

### NFR-4: Scalability
- **Horizontal Scaling**: Stateless services designed for horizontal scaling
- **Database**: PostgreSQL with connection pooling
- **Cache**: Redis for hot data (agent status, recent contexts)
- **Message Bus**: Kafka for message routing (scales horizontally)

### NFR-5: Observability
- **Metrics**: Prometheus metrics for all operations
- **Logging**: Structured JSON logs (zerolog or zap)
- **Tracing**: OpenTelemetry support (basic in Phase 1, full in Phase 2)
- **Health**: Health check endpoints

### NFR-6: Developer Experience
- **SDK**: Complete Go SDK with builder pattern
- **Documentation**: API documentation, quickstart guide
- **Examples**: Demo agents showing key features
- **Local Dev**: Docker Compose for local development

### NFR-7: Testing
- **Coverage**: >90% code coverage mandatory
- **Unit Tests**: All components have unit tests
- **Integration Tests**: Critical paths covered
- **E2E Tests**: Full agent lifecycle scenarios

---

## Architecture Overview

### Architecture Decisions (Supporting Future Phases)

#### 1. Hub-Based Architecture
- **Decision**: All communication goes through central ACB hub
- **Rationale**: Simplicity, observability, centralized security, easier multi-tenancy
- **Future**: A2A direct communication can be added in Phase 2+ without breaking changes

#### 2. Hybrid Consistency Model
- **Strong Consistency** (PostgreSQL): Agent registration, permissions, billing
- **Eventual Consistency** (Kafka): Messages, events, metrics
- **Rationale**: Right consistency model for each use case

#### 3. Three-Tier Context Handling
- **<1MB**: Direct via Kafka messages
- **1-100MB**: Streaming via gRPC with chunking
- **>100MB**: Object storage (Phase 2+)
- **Rationale**: Optimize for different data sizes

#### 4. Multi-Tenancy Ready
- **MVP**: Single tenant mode (tenant_id column exists but not enforced)
- **Phase 2**: Multi-tenancy with logical isolation
- **Design**: All tables include tenant_id, all queries filtered by tenant
- **Rationale**: No rework needed when enabling multi-tenancy

### Technology Stack

- **Language**: Go 1.21+
- **Message Bus**: Apache Kafka 3.x (Confluent client)
- **Hot Storage**: Redis 7+ with TLS
- **Warm Storage**: PostgreSQL 15+
- **Streaming**: gRPC + WebSocket
- **HTTP Router**: Gin or Chi (to be decided)
- **Security**: JWT (golang-jwt/jwt), TLS (crypto/tls)
- **Observability**: Prometheus, OpenTelemetry, zerolog/zap

### System Components

1. **Agent Registry**: Manages agent registration, discovery, status
2. **Context Manager**: Handles context CRUD, lifecycle, versioning
3. **Message Router**: Routes messages via Kafka (point-to-point, broadcast, topic-based)
4. **Streaming Service**: Handles large context transfers via gRPC
5. **Authentication Service**: JWT validation, token refresh
6. **Authorization Service**: RBAC enforcement, context-level access control
7. **HTTP Server**: REST API endpoints
8. **gRPC Server**: gRPC service endpoints
9. **Agent SDK**: Go client library for agents

---

## Complete API Specifications

### REST API Endpoints

See `api/openapi/acb-api.yaml` for complete OpenAPI 3.0 specification.

#### Authentication Endpoints
- `POST /api/v1/auth/login` - Login with credentials, receive JWT
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout and invalidate token

#### Agent Registry Endpoints
- `POST /api/v1/agents` - Register agent
- `GET /api/v1/agents/{agent_id}` - Get agent details
- `DELETE /api/v1/agents/{agent_id}` - Unregister agent
- `POST /api/v1/agents/{agent_id}/heartbeat` - Send heartbeat
- `GET /api/v1/agents` - List/discover agents (with filters)

#### Context Endpoints
- `POST /api/v1/contexts` - Create context
- `GET /api/v1/contexts/{context_id}` - Get context
- `PUT /api/v1/contexts/{context_id}` - Update context
- `DELETE /api/v1/contexts/{context_id}` - Delete context
- `GET /api/v1/contexts` - List contexts (with filters)

#### Message Endpoints
- `POST /api/v1/messages` - Send message
- `GET /api/v1/messages/{message_id}` - Get message status

#### Streaming Endpoints
- `POST /api/v1/streams/init` - Initialize stream
- `POST /api/v1/streams/{stream_id}/chunks` - Upload chunk
- `GET /api/v1/streams/{stream_id}/progress` - Get stream progress
- `GET /api/v1/streams/{stream_id}` - Download stream

#### Health & Monitoring
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics

### gRPC Services

See `api/proto/` directory for complete protobuf definitions.

#### AgentRegistry Service
```protobuf
service AgentRegistry {
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Unregister(UnregisterRequest) returns (UnregisterResponse);
  rpc Heartbeat(HeartbeatRequest) returns (HeartbeatResponse);
  rpc Discover(DiscoverRequest) returns (DiscoverResponse);
  rpc GetAgent(GetAgentRequest) returns (GetAgentResponse);
  rpc ListAgents(ListAgentsRequest) returns (ListAgentsResponse);
}
```

#### ContextService
```protobuf
service ContextService {
  rpc CreateContext(CreateContextRequest) returns (CreateContextResponse);
  rpc GetContext(GetContextRequest) returns (GetContextResponse);
  rpc UpdateContext(UpdateContextRequest) returns (UpdateContextResponse);
  rpc DeleteContext(DeleteContextRequest) returns (DeleteContextResponse);
  rpc ListContexts(ListContextsRequest) returns (ListContextsResponse);
  rpc Subscribe(SubscribeRequest) returns (stream ContextEvent);
}
```

#### StreamService
```protobuf
service StreamService {
  rpc StreamContext(stream ContextChunk) returns (StreamResponse);
  rpc ReceiveContext(ContextRequest) returns (stream ContextChunk);
  rpc GetStreamProgress(StreamProgressRequest) returns (StreamProgressResponse);
}
```

**Full API Specifications**: See `api/openapi/acb-api.yaml` (REST) and `api/proto/*.proto` (gRPC) for complete request/response schemas, field validations, and examples.

---

## Error Response Format

All API errors follow standard format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "field": "additional context"
    },
    "request_id": "correlation-id",
    "timestamp": "2025-01-XXT00:00:00Z"
  }
}
```

### Error Codes
- `UNAUTHORIZED`: Authentication failed
- `FORBIDDEN`: Authorization failed
- `NOT_FOUND`: Resource not found
- `VALIDATION_ERROR`: Request validation failed
- `RATE_LIMIT_EXCEEDED`: Rate limit exceeded
- `INTERNAL_ERROR`: Server error
- `SERVICE_UNAVAILABLE`: Service temporarily unavailable

---

## Rate Limiting

- **Default**: 1000 requests/minute per agent
- **Streaming**: 10 concurrent streams per agent
- **Bandwidth**: 100 MB/minute per agent
- Rate limit headers included in responses:
  - `X-RateLimit-Limit`: Limit per window
  - `X-RateLimit-Remaining`: Remaining in current window
  - `X-RateLimit-Reset`: Unix timestamp when limit resets

---

## Out of Scope (MVP / Phase 1)

### Explicitly Deferred to Phase 2+
- **Multi-tenancy**: Single tenant mode in MVP (architecture supports multi-tenancy)
- **Object Storage**: >100MB contexts deferred to Phase 2 (object storage integration)
- **mTLS**: Mutual TLS deferred to Phase 2
- **E2E Encryption**: End-to-end encryption deferred to Phase 2
- **Advanced Monitoring**: Full OpenTelemetry/Jaeger tracing in Phase 2
- **Schema Registry**: Schema registry integration deferred
- **Multi-Language SDKs**: Python/Java/JavaScript SDKs deferred
- **A2A Direct Communication**: Direct agent-to-agent communication deferred
- **Self-Service Signup**: Tenant provisioning deferred to Phase 3
- **Billing Integration**: Usage metering and billing deferred to Phase 3
- **Marketing Website**: Deferred to Phase 3
- **Customer Dashboard**: Deferred to Phase 3

### Out of Scope Entirely
- **Agent Marketplace**: Not planned
- **Workflow Orchestration**: Not planned
- **UI/Web Interface**: CLI and APIs only (no web UI)

---

## Success Metrics

### Phase 1 Success Metrics
- **Functionality**: 100% of MVP features working
- **Performance**: P99 latency <50ms (small messages), 100+ MB/s streaming
- **Reliability**: 99.9% uptime in dev environment
- **Test Coverage**: >90% code coverage
- **Documentation**: Complete API docs, quickstart guide
- **Developer Experience**: 2 agents can exchange messages end-to-end

### Phase 2 Success Metrics (Future)
- **Scale**: 500 concurrent agents
- **Multi-tenancy**: Logical isolation working
- **Performance**: P99 latency <50ms at scale
- **Production**: Production-ready deployment

### Phase 3 Success Metrics (Future)
- **Customers**: 10+ beta customers onboarded
- **Revenue**: Billing system operational
- **Launch**: Production launch successful

---

## Dependencies & Assumptions

### Dependencies
- Apache Kafka 3.x (message bus)
- PostgreSQL 15+ (persistent storage)
- Redis 7+ (hot storage/cache)
- Go 1.21+ (development)

### Assumptions
- Agents can connect to ACB hub over network
- Agents can authenticate with JWT tokens or API keys
- Agents implement idempotency for message handling
- Agents handle at-least-once delivery semantics

---

## Risks & Mitigations

### Risk 1: Performance at Scale
- **Risk**: MVP may not meet performance targets at 50 agents
- **Mitigation**: Load testing during development, optimize hot paths

### Risk 2: Complexity of Streaming
- **Risk**: gRPC streaming implementation may be complex
- **Mitigation**: Prototype early, use Go's native streaming support

### Risk 3: Learning Curve (Go)
- **Risk**: Developer learning Go may slow development
- **Mitigation**: Include helpful comments, progressive complexity

### Risk 4: Security Vulnerabilities
- **Risk**: Security issues in MVP
- **Mitigation**: Security review, follow Go security best practices

---

## References

- **Architecture Details**: See `CONTEXT.md` for complete architecture decisions
- **API Specifications**: See `api/openapi/acb-api.yaml` and `api/proto/*.proto`
- **Task Breakdown**: See `TASKS.md` for detailed implementation tasks
- **Development Guide**: See `DEVELOPMENT.md` (to be created)

---

**Document Owner**: Product Team  
**Last Updated**: 2025-01-XX  
**Next Review**: After Phase 1 completion


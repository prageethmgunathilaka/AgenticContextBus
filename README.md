# Agentic Context Bus (ACB)

Enterprise-grade communication framework for distributed AI agents to share context across remote locations.

## Overview

ACB enables autonomous agents (AI, RPA, automation, microservices) running in different locations to:
- Share context (state, memory, knowledge, events)
- Handle any size data: small messages to GB-scale streams
- Discover and communicate with each other
- Maintain security and observability
- Scale from prototype to enterprise production

## Quick Start

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Make (optional, but recommended)

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/prageethmgunathilaka/AgenticContextBus.git
   cd AgenticContextBus
   ```

2. **Start development environment**
   ```bash
   make dev-up
   # or
   docker-compose up -d
   ```

3. **Run database migrations**
   ```bash
   make migrate-up
   ```

4. **Start ACB server**
   ```bash
   make run-server
   # or
   go run ./cmd/acb-server
   ```

5. **Run demo agents** (in separate terminals)
   ```bash
   make run-demo-a
   make run-demo-b
   ```

## Project Structure

```
acb/
├── cmd/                    # Executable entry points
│   ├── acb-server/        # Main ACB server
│   ├── acb-cli/           # CLI tool
│   └── acb-agent-demo/    # Demo agents
├── internal/              # Private application code
│   ├── auth/              # Authentication/authorization
│   ├── registry/          # Agent registry
│   ├── context/           # Context management
│   ├── router/            # Message routing
│   ├── stream/             # Streaming handlers
│   └── storage/           # Database abstractions
├── pkg/                   # Public libraries
│   └── acb-sdk/           # Agent SDK
├── api/                   # API definitions
│   ├── proto/             # gRPC protobuf definitions
│   └── openapi/           # REST OpenAPI specs
├── migrations/            # Database migrations
├── tests/                 # Integration/E2E tests
└── docs/                  # Documentation
```

## Documentation

- **[CONTEXT.md](./CONTEXT.md)** - Complete project context and architecture decisions
- **[PRD.md](./PRD.md)** - Product Requirements Document
- **[TASKS.md](./TASKS.md)** - Implementation task breakdown
- **[API Documentation](./api/openapi/acb-api.yaml)** - OpenAPI specification

## Development

### Running Tests

```bash
# All tests with coverage
make test

# Integration tests
make test-integration

# E2E tests
make test-e2e

# Coverage report
make test-coverage
```

### Building

```bash
# Build all binaries
make build

# Format code
make fmt

# Lint code
make lint
```

## Architecture

ACB uses a hub-based architecture where all agent communication goes through a central ACB hub. This provides:
- Complete observability
- Centralized security
- Easier multi-tenancy
- Simpler implementation

See [CONTEXT.md](./CONTEXT.md) for detailed architecture decisions.

## Technology Stack

- **Language**: Go 1.21+
- **Message Bus**: Apache Kafka 3.x
- **Hot Storage**: Redis 7+ with TLS
- **Warm Storage**: PostgreSQL 15+
- **Streaming**: gRPC + WebSocket
- **HTTP Router**: Gin

## License

Apache 2.0 or Business Source License

## Status

✅ **Phase 1 (MVP) - Complete**

See [TASKS.md](./TASKS.md) for implementation progress.

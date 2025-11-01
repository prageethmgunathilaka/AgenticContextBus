# Quick Start Guide

## Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Make (optional)

## Setup

### 1. Start Development Environment

```bash
# Start all services (Kafka, Redis, PostgreSQL)
make dev-up
# or
docker-compose up -d
```

### 2. Run Database Migrations

```bash
# Connect to PostgreSQL and run migrations
psql postgres://acb:acb_password@localhost:5432/acb < migrations/001_create_agents_table.sql
psql postgres://acb:acb_password@localhost:5432/acb < migrations/002_create_contexts_table.sql
psql postgres://acb:acb_password@localhost:5432/acb < migrations/003_create_messages_table.sql
psql postgres://acb:acb_password@localhost:5432/acb < migrations/004_create_audit_log_table.sql
psql postgres://acb:acb_password@localhost:5432/acb < migrations/005_create_indexes.sql
```

### 3. Start ACB Server

```bash
# In one terminal
make run-server
# or
go run ./cmd/acb-server
```

### 4. Test the API

```bash
# Login (get JWT token)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"agent-1","password":"test"}'

# Use the access_token from response
TOKEN="<your-access-token>"

# Register an agent
curl -X POST http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "agent-1",
    "type": "ml",
    "location": "us-east-1"
  }'

# Discover agents
curl -X GET http://localhost:8080/api/v1/agents \
  -H "Authorization: Bearer $TOKEN"

# Create a context
curl -X POST http://localhost:8080/api/v1/contexts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "user-profile",
    "payload": "SGVsbG8gV29ybGQ=",
    "access_control": {
      "scope": "public"
    }
  }'

# Get context
curl -X GET http://localhost:8080/api/v1/contexts/{context_id} \
  -H "Authorization: Bearer $TOKEN"
```

## Running Demo Agents

```bash
# Terminal 1: Start Agent A
go run ./cmd/acb-agent-demo/hello-world/agent-a

# Terminal 2: Start Agent B
go run ./cmd/acb-agent-demo/hello-world/agent-b
```

## Running Tests

```bash
# Unit tests (no dependencies)
go test -v -short ./internal/models/...
go test -v -short ./internal/auth/...
go test -v -short ./internal/registry/...
go test -v -short ./internal/context/...

# All tests with coverage
make test-coverage

# Integration tests (requires Docker)
make dev-up
go test -v ./internal/storage/...
```

## Troubleshooting

### Services not starting
- Check Docker is running
- Check ports 5432, 6379, 9092 are not in use
- Check logs: `docker-compose logs`

### Database connection errors
- Ensure PostgreSQL is running: `docker-compose ps`
- Check connection string matches docker-compose.yml
- Verify migrations ran successfully

### Redis connection errors
- Ensure Redis is running
- Check password matches docker-compose.yml
- Redis password is optional for local dev

## Next Steps

- See [CONTEXT.md](./CONTEXT.md) for architecture details
- See [PRD.md](./PRD.md) for requirements
- See [TASKS.md](./TASKS.md) for implementation tasks
- See [docs/TEST_COVERAGE.md](./docs/TEST_COVERAGE.md) for test information


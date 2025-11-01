# Local Development Guide

This guide will help you set up and run the Agentic Context Bus (ACB) project locally on your machine.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Initial Setup](#initial-setup)
- [Starting the Development Environment](#starting-the-development-environment)
- [Database Migrations](#database-migrations)
- [Running the ACB Server](#running-the-acb-server)
- [Running Demo Agents](#running-demo-agents)
- [Testing](#testing)
- [Development Workflow](#development-workflow)
- [Environment Variables](#environment-variables)
- [Troubleshooting](#troubleshooting)

## Prerequisites

Before you begin, ensure you have the following installed:

### Required Software

1. **Go 1.21 or later**
   - Download from: https://golang.org/dl/
   - Verify installation:
     ```bash
     go version
     ```

2. **Docker and Docker Compose**
   - Docker Desktop: https://www.docker.com/products/docker-desktop
   - Or Docker Engine + Docker Compose
   - Verify installation:
     ```bash
     docker --version
     docker-compose --version
     ```

3. **Git**
   - Download from: https://git-scm.com/downloads
   - Verify installation:
     ```bash
     git --version
     ```

### Optional Tools

- **Make** (optional, but recommended for convenience)
  - Windows: Install via Chocolatey `choco install make`
  - macOS: `xcode-select --install`
  - Linux: `sudo apt-get install make` or `sudo yum install make`

- **PostgreSQL Client** (for manual database operations)
  - Optional, migrations can be run via scripts

- **curl** or **Postman** (for API testing)

## Initial Setup

### 1. Clone the Repository

```bash
git clone https://github.com/prageethmgunathilaka/AgenticContextBus.git
cd AgenticContextBus
```

### 2. Install Go Dependencies

```bash
go mod download
go mod tidy
```

This will download all required Go packages specified in `go.mod`.

### 3. Verify Setup

```bash
# Verify Go installation
go version

# Verify Docker is running
docker ps

# Check project structure
ls -la
```

## Starting the Development Environment

The ACB project requires several services to run locally:
- **PostgreSQL** (database)
- **Redis** (caching)
- **Kafka** (message bus)
- **Zookeeper** (required by Kafka)

### Option 1: Using Make (Recommended)

```bash
make dev-up
```

This command:
- Starts all Docker containers
- Waits for services to be ready
- Displays status messages

### Option 2: Using Docker Compose Directly

```bash
docker-compose up -d
```

The `-d` flag runs containers in detached mode (background).

### Verify Services are Running

```bash
# Check container status
docker-compose ps

# View logs
docker-compose logs

# View logs for a specific service
docker-compose logs postgres
docker-compose logs redis
docker-compose logs kafka
```

### Service Details

| Service | Port | Default Credentials | Purpose |
|---------|------|---------------------|---------|
| PostgreSQL | 5432 | User: `acb`, Password: `acb_password`, DB: `acb` | Primary database |
| Redis | 6379 | Password: `acb_redis_password` | Caching layer |
| Kafka | 9092 | None | Message bus |
| Zookeeper | 2181 | None | Kafka coordination |

### Stopping Services

```bash
# Stop all services
make dev-down
# or
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```

## Database Migrations

Before running the server, you need to set up the database schema.

### Option 1: Using Make

```bash
make migrate-up
```

### Option 2: Manual Migration

If you have `psql` installed:

```bash
# Run migrations sequentially
psql postgres://acb:acb_password@localhost:5432/acb < migrations/001_create_agents_table.sql
psql postgres://acb:acb_password@localhost:5432/acb < migrations/002_create_contexts_table.sql
psql postgres://acb:acb_password@localhost:5432/acb < migrations/003_create_messages_table.sql
psql postgres://acb:acb_password@localhost:5432/acb < migrations/004_create_audit_log_table.sql
psql postgres://acb:acb_password@localhost:5432/acb < migrations/005_create_indexes.sql
```

### Option 3: Using Go Migration Script

```bash
go run scripts/migrate.go up
```

### Verify Migrations

Connect to PostgreSQL and verify tables exist:

```bash
docker-compose exec postgres psql -U acb -d acb -c "\dt"
```

You should see:
- `agents`
- `contexts`
- `messages`
- `audit_logs`

## Running the ACB Server

### Option 1: Using Make

```bash
make run-server
```

### Option 2: Using Go Run

```bash
go run ./cmd/acb-server
```

### Option 3: Building and Running Binary

```bash
# Build the binary
make build
# or
go build -o bin/acb-server ./cmd/acb-server

# Run the binary
./bin/acb-server
# Windows:
.\bin\acb-server.exe
```

### Server Output

When the server starts successfully, you should see:

```
Starting ACB Server...
Server listening on :8080
Database connected successfully
Redis connected successfully
Kafka producer initialized
```

The server will be available at: `http://localhost:8080`

### API Endpoints

The server exposes REST API endpoints:

- **Health Check**: `GET http://localhost:8080/health`
- **API Documentation**: `GET http://localhost:8080/api/v1/docs`
- **Agent Registration**: `POST http://localhost:8080/api/v1/agents`
- **Context Management**: `POST http://localhost:8080/api/v1/contexts`

See `api/openapi/acb-api.yaml` for complete API documentation.

## Running Demo Agents

The project includes demo agents that demonstrate agent communication through ACB.

### Terminal 1: Start Agent A

```bash
make run-demo-a
# or
go run ./cmd/acb-agent-demo/hello-world/agent-a
```

### Terminal 2: Start Agent B

```bash
make run-demo-b
# or
go run ./cmd/acb-agent-demo/hello-world/agent-b
```

### What Happens

1. Agent A registers with ACB
2. Agent A discovers Agent B
3. Agent A sends a message to Agent B
4. Agent B receives and responds
5. Both agents exchange context

### Streaming Demo

For large context streaming:

```bash
go run ./cmd/acb-agent-demo/streaming-demo
```

## Testing

### Running All Tests

```bash
# Run all tests with coverage
make test

# Run tests with coverage report
make test-coverage
```

This will:
- Run all unit tests
- Generate coverage report (`coverage.out`)
- Generate HTML coverage report (`coverage.html`)
- Display coverage summary

### Running Specific Test Packages

```bash
# Unit tests (no dependencies)
go test -v ./internal/models/...
go test -v ./internal/auth/...
go test -v ./internal/registry/...
go test -v ./internal/context/...

# Storage tests (require Docker)
go test -v ./internal/storage/...

# Server tests
go test -v ./internal/server/...
```

### Running Tests with Tags

```bash
# Integration tests
make test-integration
# or
go test -v -tags=integration ./tests/integration/...

# E2E tests
make test-e2e
# or
go test -v -tags=e2e ./tests/e2e/...
```

### Test Coverage Requirements

The project requires **>90% test coverage**. Check coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

To view HTML coverage report:

```bash
go tool cover -html=coverage.out -o coverage.html
open coverage.html  # macOS
start coverage.html  # Windows
xdg-open coverage.html  # Linux
```

## Development Workflow

### 1. Make Code Changes

Edit files in your IDE or text editor.

### 2. Format Code

```bash
make fmt
# or
go fmt ./...
```

### 3. Lint Code

```bash
make lint
# or
golangci-lint run ./...
```

### 4. Run Tests

```bash
make test
```

### 5. Build

```bash
make build
```

### 6. Run Locally

```bash
make run-server
```

### Common Development Commands

```bash
# Clean build artifacts
make clean

# Format and lint
make fmt lint

# Run tests before committing
make test-coverage

# Generate protobuf code (if proto files changed)
make proto

# Install/update dependencies
make deps
```

## Environment Variables

The server can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `ACB_PORT` | `8080` | HTTP server port |
| `ACB_POSTGRES_HOST` | `localhost` | PostgreSQL host |
| `ACB_POSTGRES_PORT` | `5432` | PostgreSQL port |
| `ACB_POSTGRES_USER` | `acb` | PostgreSQL user |
| `ACB_POSTGRES_PASSWORD` | `acb_password` | PostgreSQL password |
| `ACB_POSTGRES_DB` | `acb` | PostgreSQL database name |
| `ACB_REDIS_HOST` | `localhost` | Redis host |
| `ACB_REDIS_PORT` | `6379` | Redis port |
| `ACB_REDIS_PASSWORD` | `acb_redis_password` | Redis password |
| `ACB_KAFKA_BROKERS` | `localhost:9092` | Kafka broker addresses (comma-separated) |
| `ACB_LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `ACB_JWT_SECRET` | `your-secret-key-change-in-production` | JWT signing secret |

### Setting Environment Variables

**Linux/macOS:**

```bash
export ACB_PORT=8080
export ACB_LOG_LEVEL=debug
go run ./cmd/acb-server
```

**Windows PowerShell:**

```powershell
$env:ACB_PORT=8080
$env:ACB_LOG_LEVEL=debug
go run ./cmd/acb-server
```

**Using .env file:**

Create a `.env` file in the project root:

```env
ACB_PORT=8080
ACB_LOG_LEVEL=debug
ACB_POSTGRES_HOST=localhost
```

Then load it before running:

```bash
# Linux/macOS
export $(cat .env | xargs)
go run ./cmd/acb-server

# Or use a tool like `godotenv`
go run -tags=dev ./cmd/acb-server
```

## Troubleshooting

### Services Won't Start

**Problem**: Docker containers fail to start

**Solutions**:
```bash
# Check Docker is running
docker ps

# Check port conflicts
netstat -an | grep 5432  # PostgreSQL
netstat -an | grep 6379  # Redis
netstat -an | grep 9092  # Kafka

# View container logs
docker-compose logs

# Restart services
docker-compose down
docker-compose up -d
```

### Database Connection Errors

**Problem**: Cannot connect to PostgreSQL

**Solutions**:
```bash
# Verify PostgreSQL is running
docker-compose ps postgres

# Check connection
docker-compose exec postgres psql -U acb -d acb -c "SELECT 1;"

# Verify credentials match docker-compose.yml
# Default: User: acb, Password: acb_password, DB: acb

# Reset database (WARNING: deletes all data)
docker-compose down -v
docker-compose up -d
make migrate-up
```

### Redis Connection Errors

**Problem**: Cannot connect to Redis

**Solutions**:
```bash
# Verify Redis is running
docker-compose ps redis

# Test Redis connection
docker-compose exec redis redis-cli -a acb_redis_password ping
# Should return: PONG

# Check password matches docker-compose.yml
# Default: acb_redis_password
```

### Kafka Connection Errors

**Problem**: Cannot connect to Kafka

**Solutions**:
```bash
# Verify Kafka is running
docker-compose ps kafka

# Check Kafka logs
docker-compose logs kafka

# Verify Zookeeper is running (Kafka dependency)
docker-compose ps zookeeper

# Wait for Kafka to fully start (can take 30-60 seconds)
docker-compose logs -f kafka
```

### Port Already in Use

**Problem**: Port 5432, 6379, or 9092 already in use

**Solutions**:

1. **Stop conflicting services**:
   ```bash
   # Find process using port
   # Linux/macOS
   lsof -i :5432
   # Windows
   netstat -ano | findstr :5432
   
   # Kill process or stop service
   ```

2. **Change ports in docker-compose.yml**:
   ```yaml
   ports:
     - "5433:5432"  # Change host port
   ```
   Then update environment variables accordingly.

### Go Module Errors

**Problem**: `go: cannot find module` errors

**Solutions**:
```bash
# Clear module cache
go clean -modcache

# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Verify go.mod syntax
go mod verify
```

### Test Failures

**Problem**: Tests fail with connection errors

**Solutions**:
```bash
# Ensure services are running
make dev-up

# Run tests with short flag (skip integration tests)
go test -short ./...

# Run integration tests separately
go test -tags=integration ./tests/integration/...
```

### Build Failures

**Problem**: `go build` fails

**Solutions**:
```bash
# Verify Go version (requires 1.21+)
go version

# Clean build cache
go clean -cache

# Verify all dependencies
go mod verify

# Rebuild from scratch
make clean
make build
```

### Migration Errors

**Problem**: Database migrations fail

**Solutions**:
```bash
# Check PostgreSQL is running
docker-compose ps postgres

# Verify database exists
docker-compose exec postgres psql -U acb -l

# Check migration files exist
ls migrations/

# Run migrations manually
docker-compose exec postgres psql -U acb -d acb -f /path/to/migration.sql
```

## Next Steps

- Read [CONTEXT.md](./CONTEXT.md) for architecture details
- Review [PRD.md](./PRD.md) for product requirements
- Check [TASKS.md](./TASKS.md) for implementation status
- Explore API documentation in `api/openapi/acb-api.yaml`
- See [docs/quickstart.md](./docs/quickstart.md) for quick reference

## Getting Help

- Check existing issues: https://github.com/prageethmgunathilaka/AgenticContextBus/issues
- Review logs: `docker-compose logs`
- Verify service health: `docker-compose ps`

## Contributing

When contributing code:

1. **Create a feature branch**: `git checkout -b feature/your-feature-name`
2. **Make changes**: Write code and tests
3. **Run tests**: `make test-coverage` (must be >90%)
4. **Format code**: `make fmt`
5. **Lint code**: `make lint`
6. **Commit**: `git commit -m "Description of changes"`
7. **Push**: `git push origin feature/your-feature-name`
8. **Create PR**: Open a pull request on GitHub

The CI/CD pipeline will automatically run tests and verify coverage on your PR.


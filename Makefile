.PHONY: build test clean run-server run-demo dev-up dev-down test-coverage test-integration test-e2e lint fmt proto

# Build all binaries
build:
	@echo "Building ACB server..."
	@go build -o bin/acb-server ./cmd/acb-server
	@echo "Building ACB CLI..."
	@go build -o bin/acb-cli ./cmd/acb-cli
	@echo "Building demo agents..."
	@go build -o bin/agent-a ./cmd/acb-agent-demo/hello-world/agent-a
	@go build -o bin/agent-b ./cmd/acb-agent-demo/hello-world/agent-b
	@go build -o bin/streaming-demo ./cmd/acb-agent-demo/streaming-demo

# Run tests with coverage
test:
	@echo "Running tests..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@go test -v -tags=integration ./tests/integration/...

# Run E2E tests
test-e2e:
	@echo "Running E2E tests..."
	@go test -v -tags=e2e ./tests/e2e/...

# Generate coverage report
test-coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

# Run server
run-server:
	@echo "Starting ACB server..."
	@go run ./cmd/acb-server

# Run demo agent A
run-demo-a:
	@go run ./cmd/acb-agent-demo/hello-world/agent-a

# Run demo agent B
run-demo-b:
	@go run ./cmd/acb-agent-demo/hello-world/agent-b

# Start development environment
dev-up:
	@echo "Starting development environment..."
	@docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Development environment ready!"

# Stop development environment
dev-down:
	@echo "Stopping development environment..."
	@docker-compose down

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf coverage.out coverage.html
	@go clean -cache

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run ./...

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/*.proto

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Run database migrations
migrate-up:
	@echo "Running database migrations..."
	@go run scripts/migrate.go up

# Run database migrations down
migrate-down:
	@echo "Rolling back database migrations..."
	@go run scripts/migrate.go down


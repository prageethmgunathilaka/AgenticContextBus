package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/acb/internal/auth"
	contextmgr "github.com/acb/internal/context"
	"github.com/acb/internal/registry"
	"github.com/acb/internal/server"
	"github.com/acb/internal/storage"
)

func main() {
	// Load configuration from environment
	postgresConnString := getEnv("POSTGRES_CONN_STRING", "postgres://acb:acb_password@localhost:5432/acb?sslmode=disable")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "acb_redis_password")
	httpPort := getEnv("HTTP_PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "development-secret-key-change-in-production")

	log.Println("Starting ACB Server...")
	log.Printf("HTTP Port: %s", httpPort)

	// Initialize PostgreSQL
	postgresPool, err := storage.NewPostgresStore(postgresConnString)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}
	defer postgresPool.Close()

	// Initialize Redis
	redisStore, err := storage.NewRedisStore(redisAddr, redisPassword)
	if err != nil {
		log.Printf("Warning: Failed to initialize Redis: %v (continuing without cache)", err)
		redisStore = nil
	} else {
		defer redisStore.Close()
	}

	// Initialize stores
	agentStore := storage.NewPostgresAgentStore(postgresPool.Pool())
	contextStore := storage.NewPostgresContextStore(postgresPool.Pool())

	// Initialize services
	registrySvc := registry.NewService(agentStore)
	contextMgr := contextmgr.NewManager(contextStore)

	// Initialize auth
	jwtManager := auth.NewJWTManager(jwtSecret)
	rbac := auth.NewRBAC()

	// Initialize HTTP server
	httpSrv := server.NewHTTPServer(httpPort, registrySvc, contextMgr, jwtManager, rbac)

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := httpSrv.Start(ctx); err != nil {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	log.Println("ACB Server started successfully")
	log.Printf("HTTP API available at http://localhost:%s/api/v1", httpPort)

	// Wait for shutdown signal or error
	select {
	case err := <-errChan:
		log.Fatalf("Server error: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal: %v, shutting down...", sig)
	}

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown servers gracefully
	_ = shutdownCtx

	log.Println("Server stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresStore wraps PostgreSQL connection pool
type PostgresStore struct {
	pool *pgxpool.Pool
}

// Pool returns the underlying connection pool
func (s *PostgresStore) Pool() *pgxpool.Pool {
	return s.pool
}

// NewPostgresStore creates a new PostgreSQL store
func NewPostgresStore(connString string) (*PostgresStore, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresStore{pool: pool}, nil
}

// Close closes the connection pool
func (s *PostgresStore) Close() {
	s.pool.Close()
}

// Health checks if the database is healthy
func (s *PostgresStore) Health(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

// IsUniqueViolation checks if error is a unique violation
func IsUniqueViolation(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}

// IsNotFound checks if error indicates not found
func IsNotFound(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return pgErr.Code == "02000" // no_data_found
	}
	return errors.Is(err, pgx.ErrNoRows)
}

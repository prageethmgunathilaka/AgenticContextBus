package storage

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresStore_Health(t *testing.T) {
	// Skip if no database connection
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	store, err := NewPostgresStore("postgres://acb:acb_password@localhost:5432/acb?sslmode=disable")
	require.NoError(t, err)
	defer store.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = store.Health(ctx)
	assert.NoError(t, err)
}

func TestIsUniqueViolation(t *testing.T) {
	// This test would require actual PostgreSQL error
	// For now, we test the logic
	err := &pgconn.PgError{Code: "23505"}
	assert.True(t, IsUniqueViolation(err))

	err2 := &pgconn.PgError{Code: "00000"}
	assert.False(t, IsUniqueViolation(err2))
}

func TestIsNotFound(t *testing.T) {
	// Test pgx.ErrNoRows
	assert.True(t, IsNotFound(pgx.ErrNoRows))

	// Test other error
	assert.False(t, IsNotFound(nil))
}

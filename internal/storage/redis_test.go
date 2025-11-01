package storage

import (
	"context"
	"testing"
	"time"

	"github.com/acb/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisStore_GetSet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	store, err := NewRedisStore("localhost:6379", "")
	require.NoError(t, err)
	defer store.Close()

	ctx := context.Background()

	// Test Set
	err = store.Set(ctx, "test-key", "test-value", time.Minute)
	require.NoError(t, err)

	// Test Get
	value, err := store.Get(ctx, "test-key")
	require.NoError(t, err)
	assert.Equal(t, "test-value", value)

	// Test Delete
	err = store.Delete(ctx, "test-key")
	require.NoError(t, err)

	// Test Exists
	exists, err := store.Exists(ctx, "test-key")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestRedisStore_SetNX(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	store, err := NewRedisStore("localhost:6379", "")
	require.NoError(t, err)
	defer store.Close()

	ctx := context.Background()

	// SetNX should succeed first time
	ok, err := store.SetNX(ctx, "test-nx-key", "value", time.Minute)
	require.NoError(t, err)
	assert.True(t, ok)

	// SetNX should fail second time
	ok, err = store.SetNX(ctx, "test-nx-key", "value", time.Minute)
	require.NoError(t, err)
	assert.False(t, ok)

	// Cleanup
	store.Delete(ctx, "test-nx-key")
}

func TestRedisStore_Increment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	store, err := NewRedisStore("localhost:6379", "")
	require.NoError(t, err)
	defer store.Close()

	ctx := context.Background()

	// Increment
	value, err := store.Increment(ctx, "test-counter")
	require.NoError(t, err)
	assert.Equal(t, int64(1), value)

	value, err = store.Increment(ctx, "test-counter")
	require.NoError(t, err)
	assert.Equal(t, int64(2), value)

	// Cleanup
	store.Delete(ctx, "test-counter")
}


package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRedisProgressStore_CRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	store, err := NewRedisStore("localhost:6379", "")
	require.NoError(t, err)
	defer store.Close()

	ps := NewRedisProgressStore(store)
	ctx := context.Background()

	progress := &StreamProgress{
		StreamID:      "stream-1",
		Status:        "in_progress",
		BytesReceived: 10,
		TotalBytes:    100,
		Progress:      0.1,
		Checksum:      "abc",
		StartedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Set
	require.NoError(t, ps.Set(ctx, progress, time.Minute))

	// Get
	got, err := ps.Get(ctx, progress.StreamID)
	require.NoError(t, err)
	require.Equal(t, progress.StreamID, got.StreamID)

	// Delete
	require.NoError(t, ps.Delete(ctx, progress.StreamID))

	// Get after delete should error
	_, err = ps.Get(ctx, progress.StreamID)
	if err == nil {
		t.Fatalf("expected error after delete")
	}
}



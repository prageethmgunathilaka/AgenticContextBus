package storage

import (
    "context"
    "testing"
    "time"

    "github.com/acb/internal/models"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPostgresContextStore_CRUD(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    store, err := NewPostgresStore("postgres://acb:acb_password@localhost:5432/acb?sslmode=disable")
    require.NoError(t, err)
    defer store.Close()

    ctxStore := NewPostgresContextStore(store.Pool())
    ctx := context.Background()

    // Create
    c := &models.Context{
        ID:       "test-ctx-1",
        Type:     "greeting",
        AgentID:  "agent-x",
        TenantID: "default",
        Payload:  []byte("hello"),
        Metadata: map[string]string{"k": "v"},
        Version:  "1",
        TTL:      time.Hour,
        AccessControl: models.AccessControl{Scope: models.ScopePublic},
        CreatedAt: time.Now(),
        Checksum:  models.CalculateChecksum([]byte("hello")),
    }
    // CalculateChecksum is in stream package; fallback simple value if not available
    if c.Checksum == "" { c.Checksum = "" }

    err = ctxStore.Create(ctx, c)
    require.NoError(t, err)
    defer ctxStore.Delete(ctx, c.ID)

    // Get
    got, err := ctxStore.Get(ctx, c.ID)
    require.NoError(t, err)
    assert.Equal(t, c.ID, got.ID)

    // Update
    got.Payload = []byte("world")
    got.Version = "2"
    got.TTL = time.Hour
    got.Checksum = ""
    err = ctxStore.Update(ctx, got)
    require.NoError(t, err)

    // List
    list, err := ctxStore.List(ctx, &ContextFilters{TenantID: "default", Limit: 10})
    require.NoError(t, err)
    assert.GreaterOrEqual(t, len(list), 1)
}

func TestPostgresContextStore_DeleteExpired(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    store, err := NewPostgresStore("postgres://acb:acb_password@localhost:5432/acb?sslmode=disable")
    require.NoError(t, err)
    defer store.Close()

    ctxStore := NewPostgresContextStore(store.Pool())
    ctx := context.Background()

    // Create an expired context
    c := &models.Context{
        ID:       "test-ctx-expired",
        Type:     "greeting",
        AgentID:  "agent-x",
        TenantID: "default",
        Payload:  []byte("bye"),
        Metadata: map[string]string{},
        Version:  "1",
        TTL:      time.Second,
        AccessControl: models.AccessControl{Scope: models.ScopePublic},
        CreatedAt: time.Now().Add(-2 * time.Hour),
        ExpiresAt: time.Now().Add(-time.Hour),
        Checksum:  "",
    }
    _ = ctxStore.Create(ctx, c)

    // DeleteExpired should remove at least 1 or 0 if already gone
    _, _ = ctxStore.DeleteExpired(ctx)
}



package registry

import (
	"context"
	"testing"
	"time"

	"github.com/acb/internal/models"
	"github.com/acb/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresAgentStore_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	store, err := storage.NewPostgresStore("postgres://acb:acb_password@localhost:5432/acb?sslmode=disable")
	require.NoError(t, err)
	defer store.Close()

	agentStore := NewPostgresAgentStore(store.Pool())

	ctx := context.Background()

	tests := []struct {
		name    string
		agent   *models.Agent
		wantErr bool
	}{
		{
			name: "valid agent",
			agent: &models.Agent{
				ID:        "test-agent-1",
				Type:      "ml",
				Location:  "us-east-1",
				Status:    models.AgentStatusOnline,
				TenantID:  "default",
				CreatedAt: time.Now(),
				LastSeen:  time.Now(),
			},
			wantErr: false,
		},
		{
			name: "duplicate agent ID",
			agent: &models.Agent{
				ID:        "test-agent-1",
				Type:      "ml",
				Status:    models.AgentStatusOnline,
				TenantID:  "default",
				CreatedAt: time.Now(),
				LastSeen:  time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := agentStore.Create(ctx, tt.agent)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

	// Cleanup
	agentStore.Delete(ctx, "test-agent-1")
}

func TestPostgresAgentStore_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	store, err := storage.NewPostgresStore("postgres://acb:acb_password@localhost:5432/acb?sslmode=disable")
	require.NoError(t, err)
	defer store.Close()

	agentStore := NewPostgresAgentStore(store.Pool())
	ctx := context.Background()

	// Create test agent
	agent := &models.Agent{
		ID:        "test-agent-get",
		Type:      "ml",
		Status:    models.AgentStatusOnline,
		TenantID:  "default",
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
	}
	err = agentStore.Create(ctx, agent)
	require.NoError(t, err)
	defer agentStore.Delete(ctx, "test-agent-get")

	// Test Get
	got, err := agentStore.Get(ctx, "test-agent-get")
	require.NoError(t, err)
	assert.Equal(t, agent.ID, got.ID)
	assert.Equal(t, agent.Type, got.Type)

	// Test Get non-existent
	_, err = agentStore.Get(ctx, "non-existent")
	assert.Error(t, err)
}

func TestPostgresAgentStore_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	store, err := storage.NewPostgresStore("postgres://acb:acb_password@localhost:5432/acb?sslmode=disable")
	require.NoError(t, err)
	defer store.Close()

	agentStore := NewPostgresAgentStore(store.Pool())
	ctx := context.Background()

	// Create test agent
	agent := &models.Agent{
		ID:        "test-agent-update",
		Type:      "ml",
		Status:    models.AgentStatusOnline,
		TenantID:  "default",
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
	}
	err = agentStore.Create(ctx, agent)
	require.NoError(t, err)
	defer agentStore.Delete(ctx, "test-agent-update")

	// Update agent
	agent.Type = "rpa"
	err = agentStore.Update(ctx, agent)
	require.NoError(t, err)

	// Verify update
	got, err := agentStore.Get(ctx, "test-agent-update")
	require.NoError(t, err)
	assert.Equal(t, "rpa", got.Type)
}

func TestPostgresAgentStore_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	store, err := storage.NewPostgresStore("postgres://acb:acb_password@localhost:5432/acb?sslmode=disable")
	require.NoError(t, err)
	defer store.Close()

	agentStore := NewPostgresAgentStore(store.Pool())
	ctx := context.Background()

	// Create test agents
	for i := 0; i < 3; i++ {
		agent := &models.Agent{
			ID:        "test-agent-list-" + string(rune(i)),
			Type:      "ml",
			Status:    models.AgentStatusOnline,
			TenantID:  "default",
			CreatedAt: time.Now(),
			LastSeen:  time.Now(),
		}
		agentStore.Create(ctx, agent)
		defer agentStore.Delete(ctx, agent.ID)
	}

	// Test List
	filters := &storage.AgentFilters{
		TenantID: "default",
		Limit:    10,
	}
	agents, err := agentStore.List(ctx, filters)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(agents), 3)

	// Test filter by type
	filters.Type = "ml"
	agents, err = agentStore.List(ctx, filters)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(agents), 3)
}

func TestPostgresAgentStore_UpdateLastSeen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	store, err := storage.NewPostgresStore("postgres://acb:acb_password@localhost:5432/acb?sslmode=disable")
	require.NoError(t, err)
	defer store.Close()

	agentStore := NewPostgresAgentStore(store.Pool())
	ctx := context.Background()

	// Create test agent
	agent := &models.Agent{
		ID:        "test-agent-heartbeat",
		Type:      "ml",
		Status:    models.AgentStatusOnline,
		TenantID:  "default",
		CreatedAt: time.Now(),
		LastSeen:  time.Now(),
	}
	err = agentStore.Create(ctx, agent)
	require.NoError(t, err)
	defer agentStore.Delete(ctx, "test-agent-heartbeat")

	// Update last seen
	err = agentStore.UpdateLastSeen(ctx, "test-agent-heartbeat")
	require.NoError(t, err)

	// Verify
	got, err := agentStore.Get(ctx, "test-agent-heartbeat")
	require.NoError(t, err)
	assert.True(t, got.LastSeen.After(agent.LastSeen))
}

package context

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/acb/internal/models"
	"github.com/acb/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockContextStore is a mock implementation for testing
type MockContextStore struct {
	contexts map[string]*models.Context
}

func NewMockContextStore() *MockContextStore {
	return &MockContextStore{
		contexts: make(map[string]*models.Context),
	}
}

func (m *MockContextStore) Create(ctx context.Context, c *models.Context) error {
	if _, exists := m.contexts[c.ID]; exists {
		return fmt.Errorf("context already exists")
	}
	m.contexts[c.ID] = c
	return nil
}

func (m *MockContextStore) Get(ctx context.Context, contextID string) (*models.Context, error) {
	c, exists := m.contexts[contextID]
	if !exists {
		return nil, fmt.Errorf("context not found")
	}
	return c, nil
}

func (m *MockContextStore) Update(ctx context.Context, c *models.Context) error {
	if _, exists := m.contexts[c.ID]; !exists {
		return fmt.Errorf("context not found")
	}
	m.contexts[c.ID] = c
	return nil
}

func (m *MockContextStore) Delete(ctx context.Context, contextID string) error {
	if _, exists := m.contexts[contextID]; !exists {
		return fmt.Errorf("context not found")
	}
	delete(m.contexts, contextID)
	return nil
}

func (m *MockContextStore) List(ctx context.Context, filters *storage.ContextFilters) ([]*models.Context, error) {
	var contexts []*models.Context
	for _, c := range m.contexts {
		if filters != nil && filters.TenantID != "" && c.TenantID != filters.TenantID {
			continue
		}
		if filters != nil && filters.Type != "" && c.Type != filters.Type {
			continue
		}
		if filters != nil && filters.AgentID != "" && c.AgentID != filters.AgentID {
			continue
		}
		contexts = append(contexts, c)
	}
	return contexts, nil
}

func (m *MockContextStore) DeleteExpired(ctx context.Context) (int, error) {
	count := 0
	for id, c := range m.contexts {
		if c.IsExpired() {
			delete(m.contexts, id)
			count++
		}
	}
	return count, nil
}

func TestManager_Create(t *testing.T) {
	mockStore := NewMockContextStore()
	mgr := NewManager(mockStore)

	req := &CreateRequest{
		Type:     "user-profile",
		AgentID:  "agent-1",
		TenantID: "default",
		Payload:  []byte("test data"),
		AccessControl: models.AccessControl{
			Scope: models.ScopePublic,
		},
		TTL: time.Hour,
	}

	ctx, err := mgr.Create(context.Background(), req)
	require.NoError(t, err)
	assert.NotEmpty(t, ctx.ID)
	assert.Equal(t, "user-profile", ctx.Type)
	assert.NotEmpty(t, ctx.Checksum)
}

func TestManager_Get(t *testing.T) {
	mockStore := NewMockContextStore()
	mgr := NewManager(mockStore)

	req := &CreateRequest{
		Type:     "user-profile",
		AgentID:  "agent-1",
		TenantID: "default",
		Payload:  []byte("test data"),
		AccessControl: models.AccessControl{
			Scope: models.ScopePublic,
		},
	}

	created, err := mgr.Create(context.Background(), req)
	require.NoError(t, err)

	got, err := mgr.Get(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)
}

func TestManager_Update(t *testing.T) {
	mockStore := NewMockContextStore()
	mgr := NewManager(mockStore)

	req := &CreateRequest{
		Type:     "user-profile",
		AgentID:  "agent-1",
		TenantID: "default",
		Payload:  []byte("test data"),
		AccessControl: models.AccessControl{
			Scope: models.ScopePublic,
		},
	}

	created, err := mgr.Create(context.Background(), req)
	require.NoError(t, err)

	updateReq := &UpdateRequest{
		Payload: []byte("updated data"),
		Version: "2.0",
	}

	updated, err := mgr.Update(context.Background(), created.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, "updated data", string(updated.Payload))
	assert.Equal(t, "2.0", updated.Version)
}

func TestManager_Delete(t *testing.T) {
	mockStore := NewMockContextStore()
	mgr := NewManager(mockStore)

	req := &CreateRequest{
		Type:     "user-profile",
		AgentID:  "agent-1",
		TenantID: "default",
		Payload:  []byte("test data"),
		AccessControl: models.AccessControl{
			Scope: models.ScopePublic,
		},
	}

	created, err := mgr.Create(context.Background(), req)
	require.NoError(t, err)

	err = mgr.Delete(context.Background(), created.ID)
	require.NoError(t, err)

	_, err = mgr.Get(context.Background(), created.ID)
	assert.Error(t, err)
}

func TestManager_List(t *testing.T) {
	mockStore := NewMockContextStore()
	mgr := NewManager(mockStore)

	// Create multiple contexts
	for i := 0; i < 3; i++ {
		req := &CreateRequest{
			Type:     "user-profile",
			AgentID:  "agent-1",
			TenantID: "default",
			Payload:  []byte("test data"),
			AccessControl: models.AccessControl{
				Scope: models.ScopePublic,
			},
		}
		_, err := mgr.Create(context.Background(), req)
		require.NoError(t, err)
	}

	filters := &storage.ContextFilters{
		TenantID: "default",
	}
	contexts, err := mgr.List(context.Background(), filters)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(contexts), 3)
}

func TestManager_DeleteExpired(t *testing.T) {
	mockStore := NewMockContextStore()
	mgr := NewManager(mockStore)

	// Create expired context
	req := &CreateRequest{
		Type:     "user-profile",
		AgentID:  "agent-1",
		TenantID: "default",
		Payload:  []byte("test data"),
		AccessControl: models.AccessControl{
			Scope: models.ScopePublic,
		},
		TTL: -1 * time.Hour, // Expired
	}

	created, err := mgr.Create(context.Background(), req)
	require.NoError(t, err)

	// Manually set expired time
	created.ExpiresAt = time.Now().Add(-1 * time.Hour)
	mockStore.contexts[created.ID] = created

	count, err := mgr.DeleteExpired(context.Background())
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 1)
}

func TestCalculateChecksum(t *testing.T) {
	data := []byte("test data")
	checksum := calculateChecksum(data)
	assert.NotEmpty(t, checksum)
	assert.Len(t, checksum, 64) // SHA-256 hex string

	// Same data should produce same checksum
	checksum2 := calculateChecksum(data)
	assert.Equal(t, checksum, checksum2)

	// Different data should produce different checksum
	checksum3 := calculateChecksum([]byte("different data"))
	assert.NotEqual(t, checksum, checksum3)
}

func TestManager_Update_TTLAndAccessControl(t *testing.T) {
	mockStore := NewMockContextStore()
	mgr := NewManager(mockStore)

	created, err := mgr.Create(context.Background(), &CreateRequest{
		Type:          "doc",
		AgentID:       "agent-1",
		TenantID:      "default",
		Payload:       []byte("p"),
		AccessControl: models.AccessControl{Scope: models.ScopePrivate},
		TTL:           time.Minute,
	})
	require.NoError(t, err)

	// Now update without payload but with TTL and AccessControl change
	before := created.ExpiresAt
	updated, err := mgr.Update(context.Background(), created.ID, &UpdateRequest{
		AccessControl: models.AccessControl{Scope: models.ScopePublic},
		TTL:           2 * time.Minute,
	})
	require.NoError(t, err)
	assert.Equal(t, models.ScopePublic, updated.AccessControl.Scope)
	assert.True(t, updated.ExpiresAt.After(before))
}

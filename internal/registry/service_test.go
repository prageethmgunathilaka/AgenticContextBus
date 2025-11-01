package registry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/acb/internal/models"
	"github.com/acb/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ErrDuplicate = errors.New("duplicate")
	ErrNotFound  = errors.New("not found")
)

// MockAgentStore is a mock implementation for testing
type MockAgentStore struct {
	agents map[string]*models.Agent
}

func NewMockAgentStore() *MockAgentStore {
	return &MockAgentStore{
		agents: make(map[string]*models.Agent),
	}
}

func (m *MockAgentStore) Create(ctx context.Context, agent *models.Agent) error {
	if _, exists := m.agents[agent.ID]; exists {
		return fmt.Errorf("agent with ID %s already exists: %w", agent.ID, ErrDuplicate)
	}
	m.agents[agent.ID] = agent
	return nil
}

func (m *MockAgentStore) Get(ctx context.Context, agentID string) (*models.Agent, error) {
	agent, exists := m.agents[agentID]
	if !exists {
		return nil, fmt.Errorf("agent not found: %w", ErrNotFound)
	}
	return agent, nil
}

func (m *MockAgentStore) Update(ctx context.Context, agent *models.Agent) error {
	if _, exists := m.agents[agent.ID]; !exists {
		return fmt.Errorf("agent not found: %w", ErrNotFound)
	}
	m.agents[agent.ID] = agent
	return nil
}

func (m *MockAgentStore) Delete(ctx context.Context, agentID string) error {
	if _, exists := m.agents[agentID]; !exists {
		return fmt.Errorf("agent not found: %w", ErrNotFound)
	}
	delete(m.agents, agentID)
	return nil
}

func (m *MockAgentStore) List(ctx context.Context, filters *storage.AgentFilters) ([]*models.Agent, error) {
	var agents []*models.Agent
	for _, agent := range m.agents {
		if filters != nil && filters.TenantID != "" && agent.TenantID != filters.TenantID {
			continue
		}
		if filters != nil && filters.Type != "" && agent.Type != filters.Type {
			continue
		}
		agents = append(agents, agent)
	}
	return agents, nil
}

func (m *MockAgentStore) UpdateLastSeen(ctx context.Context, agentID string) error {
	agent, exists := m.agents[agentID]
	if !exists {
		return fmt.Errorf("agent not found: %w", ErrNotFound)
	}
	agent.LastSeen = time.Now()
	agent.Status = models.AgentStatusOnline
	return nil
}

func TestService_Register(t *testing.T) {
	mockStore := NewMockAgentStore()
	service := NewService(mockStore)

	req := &RegisterRequest{
		ID:       "test-agent",
		Type:     "ml",
		TenantID: "default",
	}

	agent, err := service.Register(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "test-agent", agent.ID)
	assert.Equal(t, "ml", agent.Type)
}

func TestService_Register_Duplicate(t *testing.T) {
	mockStore := NewMockAgentStore()
	service := NewService(mockStore)

	req := &RegisterRequest{
		ID:       "test-agent",
		Type:     "ml",
		TenantID: "default",
	}

	_, err := service.Register(context.Background(), req)
	require.NoError(t, err)

	// Try to register again
	_, err = service.Register(context.Background(), req)
	assert.Error(t, err)
}

func TestService_Get(t *testing.T) {
	mockStore := NewMockAgentStore()
	service := NewService(mockStore)

	req := &RegisterRequest{
		ID:       "test-agent",
		Type:     "ml",
		TenantID: "default",
	}

	_, err := service.Register(context.Background(), req)
	require.NoError(t, err)

	agent, err := service.Get(context.Background(), "test-agent")
	require.NoError(t, err)
	assert.Equal(t, "test-agent", agent.ID)
}

func TestService_Unregister(t *testing.T) {
	mockStore := NewMockAgentStore()
	service := NewService(mockStore)

	req := &RegisterRequest{
		ID:       "test-agent",
		Type:     "ml",
		TenantID: "default",
	}

	_, err := service.Register(context.Background(), req)
	require.NoError(t, err)

	err = service.Unregister(context.Background(), "test-agent")
	require.NoError(t, err)

	_, err = service.Get(context.Background(), "test-agent")
	assert.Error(t, err)
}

func TestService_Heartbeat(t *testing.T) {
	mockStore := NewMockAgentStore()
	service := NewService(mockStore)

	req := &RegisterRequest{
		ID:       "test-agent",
		Type:     "ml",
		TenantID: "default",
	}

	_, err := service.Register(context.Background(), req)
	require.NoError(t, err)

	err = service.Heartbeat(context.Background(), "test-agent")
	require.NoError(t, err)
}

func TestService_Discover(t *testing.T) {
	mockStore := NewMockAgentStore()
	service := NewService(mockStore)

	// Register multiple agents
	for i := 0; i < 3; i++ {
		req := &RegisterRequest{
			ID:       "test-agent-" + string(rune(i)),
			Type:     "ml",
			TenantID: "default",
		}
		_, err := service.Register(context.Background(), req)
		require.NoError(t, err)
	}

	filters := &storage.AgentFilters{
		TenantID: "default",
	}
	agents, err := service.Discover(context.Background(), filters)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(agents), 3)
}

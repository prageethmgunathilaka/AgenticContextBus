package registry

import (
	"context"
	"fmt"
	"time"

	"github.com/acb/internal/models"
	"github.com/acb/internal/storage"
)

// Service provides agent registry operations
type Service struct {
	store storage.AgentStore
	// cache storage.AgentCache // Will be implemented later
}

// NewService creates a new registry service
func NewService(store storage.AgentStore) *Service {
	return &Service{
		store: store,
	}
}

// Register registers a new agent
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*models.Agent, error) {
	agent := &models.Agent{
		ID:           req.ID,
		Type:         req.Type,
		Location:     req.Location,
		Capabilities: req.Capabilities,
		Metadata:     req.Metadata,
		Status:       models.AgentStatusUnknown,
		TenantID:     req.TenantID,
		CreatedAt:    time.Now(),
		LastSeen:     time.Now(),
	}

	if err := agent.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.store.Create(ctx, agent); err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	agent.Status = models.AgentStatusOnline
	return agent, nil
}

// Get retrieves an agent by ID
func (s *Service) Get(ctx context.Context, agentID string) (*models.Agent, error) {
	return s.store.Get(ctx, agentID)
}

// Unregister removes an agent
func (s *Service) Unregister(ctx context.Context, agentID string) error {
	return s.store.Delete(ctx, agentID)
}

// Heartbeat updates agent's last seen timestamp
func (s *Service) Heartbeat(ctx context.Context, agentID string) error {
	return s.store.UpdateLastSeen(ctx, agentID)
}

// Discover finds agents matching filters
func (s *Service) Discover(ctx context.Context, filters *storage.AgentFilters) ([]*models.Agent, error) {
	return s.store.List(ctx, filters)
}

// RegisterRequest contains agent registration data
type RegisterRequest struct {
	ID           string
	Type         string
	Location     string
	Capabilities []string
	Metadata     map[string]string
	TenantID     string
}

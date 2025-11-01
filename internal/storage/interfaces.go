package storage

import (
	"context"
	"time"

	"github.com/acb/internal/models"
)

// AgentStore defines the interface for agent storage operations
type AgentStore interface {
	Create(ctx context.Context, agent *models.Agent) error
	Get(ctx context.Context, agentID string) (*models.Agent, error)
	Update(ctx context.Context, agent *models.Agent) error
	Delete(ctx context.Context, agentID string) error
	List(ctx context.Context, filters *AgentFilters) ([]*models.Agent, error)
	UpdateLastSeen(ctx context.Context, agentID string) error
}

// AgentFilters contains filters for listing agents
type AgentFilters struct {
	Type       string
	Location   string
	Capability string
	Status     models.AgentStatus
	TenantID   string
	Limit      int
	Offset     int
}

// ContextStore defines the interface for context storage operations
type ContextStore interface {
	Create(ctx context.Context, c *models.Context) error
	Get(ctx context.Context, contextID string) (*models.Context, error)
	Update(ctx context.Context, c *models.Context) error
	Delete(ctx context.Context, contextID string) error
	List(ctx context.Context, filters *ContextFilters) ([]*models.Context, error)
	DeleteExpired(ctx context.Context) (int, error)
}

// ContextFilters contains filters for listing contexts
type ContextFilters struct {
	Type     string
	AgentID  string
	TenantID string
	Limit    int
	Offset   int
}

// AgentCache defines interface for agent caching
type AgentCache interface {
	Get(ctx context.Context, agentID string) (*models.Agent, error)
	Set(ctx context.Context, agent *models.Agent, ttl time.Duration) error
	Delete(ctx context.Context, agentID string) error
}

// ContextCache defines interface for context caching
type ContextCache interface {
	Get(ctx context.Context, contextID string) (*models.Context, error)
	Set(ctx context.Context, c *models.Context, ttl time.Duration) error
	Delete(ctx context.Context, contextID string) error
}


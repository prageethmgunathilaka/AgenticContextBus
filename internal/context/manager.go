package context

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/acb/internal/constants"
	"github.com/acb/internal/models"
	"github.com/acb/internal/storage"
	"github.com/google/uuid"
)

// Manager provides context management operations
type Manager struct {
	store storage.ContextStore
	cache storage.ContextCache // Will be implemented
}

// NewManager creates a new context manager
func NewManager(store storage.ContextStore) *Manager {
	return &Manager{
		store: store,
	}
}

// Create creates a new context
func (m *Manager) Create(ctx context.Context, req *CreateRequest) (*models.Context, error) {
	// Generate context ID if not provided
	contextID := req.ID
	if contextID == "" {
		contextID = uuid.New().String()
	}

	// Calculate checksum
	checksum := calculateChecksum(req.Payload)

	// Calculate expiration
	ttl := req.TTL
	if ttl == 0 {
		ttl = time.Duration(constants.DefaultContextTTL) * time.Second
	}
	expiresAt := time.Now().Add(ttl)

	c := &models.Context{
		ID:            contextID,
		Type:          req.Type,
		AgentID:       req.AgentID,
		TenantID:      req.TenantID,
		Payload:       req.Payload,
		Metadata:      req.Metadata,
		Version:       req.Version,
		AccessControl: req.AccessControl,
		TTL:           ttl,
		ExpiresAt:     expiresAt,
		Checksum:      checksum,
		CreatedAt:     time.Now(),
	}

	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := m.store.Create(ctx, c); err != nil {
		return nil, fmt.Errorf("failed to create context: %w", err)
	}

	return c, nil
}

// Get retrieves a context by ID
func (m *Manager) Get(ctx context.Context, contextID string) (*models.Context, error) {
	return m.store.Get(ctx, contextID)
}

// Update updates an existing context
func (m *Manager) Update(ctx context.Context, contextID string, req *UpdateRequest) (*models.Context, error) {
	c, err := m.store.Get(ctx, contextID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if len(req.Payload) > 0 {
		c.Payload = req.Payload
		c.Checksum = calculateChecksum(req.Payload)
	}
	if req.Metadata != nil {
		c.Metadata = req.Metadata
	}
	if req.Version != "" {
		c.Version = req.Version
	}
	if req.AccessControl.Scope != "" {
		c.AccessControl = req.AccessControl
	}
	if req.TTL > 0 {
		c.TTL = req.TTL
		c.ExpiresAt = time.Now().Add(req.TTL)
	}

	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := m.store.Update(ctx, c); err != nil {
		return nil, fmt.Errorf("failed to update context: %w", err)
	}

	return c, nil
}

// Delete removes a context
func (m *Manager) Delete(ctx context.Context, contextID string) error {
	return m.store.Delete(ctx, contextID)
}

// List retrieves contexts with filters
func (m *Manager) List(ctx context.Context, filters *storage.ContextFilters) ([]*models.Context, error) {
	return m.store.List(ctx, filters)
}

// DeleteExpired removes expired contexts
func (m *Manager) DeleteExpired(ctx context.Context) (int, error) {
	return m.store.DeleteExpired(ctx)
}

// calculateChecksum calculates SHA-256 checksum
func calculateChecksum(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// CreateRequest contains context creation data
type CreateRequest struct {
	ID            string
	Type          string
	AgentID       string
	TenantID      string
	Payload       []byte
	Metadata      map[string]string
	Version       string
	AccessControl models.AccessControl
	TTL           time.Duration
}

// UpdateRequest contains context update data
type UpdateRequest struct {
	Payload       []byte
	Metadata      map[string]string
	Version       string
	AccessControl models.AccessControl
	TTL           time.Duration
}


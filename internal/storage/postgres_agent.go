package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/acb/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresAgentStore implements AgentStore using PostgreSQL
type PostgresAgentStore struct {
	pool *pgxpool.Pool
}

// NewPostgresAgentStore creates a new PostgreSQL agent store
func NewPostgresAgentStore(pool *pgxpool.Pool) *PostgresAgentStore {
	return &PostgresAgentStore{pool: pool}
}

// Create inserts a new agent
func (s *PostgresAgentStore) Create(ctx context.Context, agent *models.Agent) error {
	if err := agent.Validate(); err != nil {
		return err
	}

    // Handle nil capabilities - pass nil instead of marshaling to avoid "null" string
    var capabilities interface{}
    if len(agent.Capabilities) > 0 {
        capabilities = agent.Capabilities
    }

	metadataJSON, _ := json.Marshal(agent.Metadata)

	query := `
		INSERT INTO agents (id, type, location, capabilities, metadata, status, tenant_id, created_at, last_seen)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := s.pool.Exec(ctx, query,
		agent.ID,
		agent.Type,
		agent.Location,
		capabilities,
		metadataJSON,
		string(agent.Status),
		agent.TenantID,
		agent.CreatedAt,
		agent.LastSeen,
	)

	if IsUniqueViolation(err) {
		return fmt.Errorf("agent with ID %s already exists", agent.ID)
	}
	return err
}

// Get retrieves an agent by ID
func (s *PostgresAgentStore) Get(ctx context.Context, agentID string) (*models.Agent, error) {
	query := `
		SELECT id, type, location, capabilities, metadata, status, tenant_id, created_at, last_seen
		FROM agents
		WHERE id = $1
	`

	var agent models.Agent
	var capabilities []string
	var metadataJSON []byte
	var statusStr string

	err := s.pool.QueryRow(ctx, query, agentID).Scan(
		&agent.ID,
		&agent.Type,
		&agent.Location,
		&capabilities,
		&metadataJSON,
		&statusStr,
		&agent.TenantID,
		&agent.CreatedAt,
		&agent.LastSeen,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("agent not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	agent.Status = models.AgentStatus(statusStr)
	agent.Capabilities = capabilities
	_ = json.Unmarshal(metadataJSON, &agent.Metadata)

	return &agent, nil
}

// Update updates an agent
func (s *PostgresAgentStore) Update(ctx context.Context, agent *models.Agent) error {
	if err := agent.Validate(); err != nil {
		return err
	}

    // Handle nil capabilities - pass nil instead of marshaling to avoid "null" string
    var capabilities interface{}
    if len(agent.Capabilities) > 0 {
        capabilities = agent.Capabilities
    }

	metadataJSON, _ := json.Marshal(agent.Metadata)

	query := `
		UPDATE agents
		SET type = $2, location = $3, capabilities = $4, metadata = $5, status = $6, last_seen = $7
		WHERE id = $1
	`

	result, err := s.pool.Exec(ctx, query,
		agent.ID,
		agent.Type,
		agent.Location,
		capabilities,
		metadataJSON,
		string(agent.Status),
		agent.LastSeen,
	)

	if result.RowsAffected() == 0 {
		return fmt.Errorf("agent not found")
	}
	return err
}

// Delete removes an agent
func (s *PostgresAgentStore) Delete(ctx context.Context, agentID string) error {
	query := `DELETE FROM agents WHERE id = $1`
	result, err := s.pool.Exec(ctx, query, agentID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("agent not found")
	}
	return nil
}

// List retrieves agents with filters
func (s *PostgresAgentStore) List(ctx context.Context, filters *AgentFilters) ([]*models.Agent, error) {
	query := `SELECT id, type, location, capabilities, metadata, status, tenant_id, created_at, last_seen FROM agents WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filters == nil {
		filters = &AgentFilters{Limit: 100}
	}

	if filters.TenantID != "" {
		query += fmt.Sprintf(" AND tenant_id = $%d", argIndex)
		args = append(args, filters.TenantID)
		argIndex++
	}

	if filters.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, filters.Type)
		argIndex++
	}

	if filters.Location != "" {
		query += fmt.Sprintf(" AND location = $%d", argIndex)
		args = append(args, filters.Location)
		argIndex++
	}

	if filters.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, string(filters.Status))
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filters.Offset)
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*models.Agent
	for rows.Next() {
		var agent models.Agent
		var capabilities []string
		var metadataJSON []byte
		var statusStr string

		err := rows.Scan(
			&agent.ID,
			&agent.Type,
			&agent.Location,
			&capabilities,
			&metadataJSON,
			&statusStr,
			&agent.TenantID,
			&agent.CreatedAt,
			&agent.LastSeen,
		)
		if err != nil {
			return nil, err
		}

		agent.Status = models.AgentStatus(statusStr)
		agent.Capabilities = capabilities
		_ = json.Unmarshal(metadataJSON, &agent.Metadata)
		agents = append(agents, &agent)
	}

	return agents, rows.Err()
}

// UpdateLastSeen updates the agent's last seen timestamp
func (s *PostgresAgentStore) UpdateLastSeen(ctx context.Context, agentID string) error {
	query := `
		UPDATE agents
		SET last_seen = $1, status = 'online'
		WHERE id = $2
	`
	result, err := s.pool.Exec(ctx, query, time.Now(), agentID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("agent not found")
	}
	return nil
}

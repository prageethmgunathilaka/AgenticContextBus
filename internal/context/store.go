package context

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/acb/internal/models"
	"github.com/acb/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresContextStore implements ContextStore using PostgreSQL
type PostgresContextStore struct {
	pool *pgxpool.Pool
}

// NewPostgresContextStore creates a new PostgreSQL context store
func NewPostgresContextStore(pool *pgxpool.Pool) *PostgresContextStore {
	return &PostgresContextStore{pool: pool}
}

// Create inserts a new context
func (s *PostgresContextStore) Create(ctx context.Context, c *models.Context) error {
	if err := c.Validate(); err != nil {
		return err
	}

	metadataJSON, _ := json.Marshal(c.Metadata)
	accessControlJSON, _ := json.Marshal(c.AccessControl)
	payloadRefJSON, _ := json.Marshal(c.PayloadRef)

	query := `
		INSERT INTO contexts (id, type, agent_id, tenant_id, payload, payload_ref, metadata, version, schema_id, ttl_seconds, access_control, created_at, expires_at, checksum)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	ttlSeconds := int(c.TTL.Seconds())
	if ttlSeconds == 0 {
		ttlSeconds = 86400 // Default 24 hours
	}

	expiresAt := c.ExpiresAt
	if expiresAt.IsZero() && c.TTL > 0 {
		expiresAt = c.CalculateExpiration()
	}

	_, err := s.pool.Exec(ctx, query,
		c.ID,
		c.Type,
		c.AgentID,
		c.TenantID,
		c.Payload,
		payloadRefJSON,
		metadataJSON,
		c.Version,
		c.SchemaID,
		ttlSeconds,
		accessControlJSON,
		c.CreatedAt,
		expiresAt,
		c.Checksum,
	)

	if storage.IsUniqueViolation(err) {
		return fmt.Errorf("context with ID %s already exists", c.ID)
	}
	return err
}

// Get retrieves a context by ID
func (s *PostgresContextStore) Get(ctx context.Context, contextID string) (*models.Context, error) {
	query := `
		SELECT id, type, agent_id, tenant_id, payload, payload_ref, metadata, version, schema_id, ttl_seconds, access_control, created_at, expires_at, checksum
		FROM contexts
		WHERE id = $1
	`

	var c models.Context
	var payloadRefJSON []byte
	var metadataJSON []byte
	var accessControlJSON []byte
	var ttlSeconds int
	var expiresAt sql.NullTime

	err := s.pool.QueryRow(ctx, query, contextID).Scan(
		&c.ID,
		&c.Type,
		&c.AgentID,
		&c.TenantID,
		&c.Payload,
		&payloadRefJSON,
		&metadataJSON,
		&c.Version,
		&c.SchemaID,
		&ttlSeconds,
		&accessControlJSON,
		&c.CreatedAt,
		&expiresAt,
		&c.Checksum,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("context not found: %w", err)
	}
	if err != nil {
		return nil, err
	}

	c.TTL = time.Duration(ttlSeconds) * time.Second
	if expiresAt.Valid {
		c.ExpiresAt = expiresAt.Time
	}

	if err := json.Unmarshal(metadataJSON, &c.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}
	if err := json.Unmarshal(accessControlJSON, &c.AccessControl); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access_control: %w", err)
	}
	if len(payloadRefJSON) > 0 {
		if err := json.Unmarshal(payloadRefJSON, &c.PayloadRef); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payload_ref: %w", err)
		}
	}

	return &c, nil
}

// Update updates a context
func (s *PostgresContextStore) Update(ctx context.Context, c *models.Context) error {
	if err := c.Validate(); err != nil {
		return err
	}

	metadataJSON, _ := json.Marshal(c.Metadata)
	accessControlJSON, _ := json.Marshal(c.AccessControl)
	payloadRefJSON, _ := json.Marshal(c.PayloadRef)

	ttlSeconds := int(c.TTL.Seconds())
	if ttlSeconds == 0 {
		ttlSeconds = 86400
	}

	expiresAt := c.ExpiresAt
	if expiresAt.IsZero() && c.TTL > 0 {
		expiresAt = c.CalculateExpiration()
	}

	query := `
		UPDATE contexts
		SET payload = $2, payload_ref = $3, metadata = $4, version = $5, ttl_seconds = $6, access_control = $7, expires_at = $8, checksum = $9
		WHERE id = $1
	`

	result, err := s.pool.Exec(ctx, query,
		c.ID,
		c.Payload,
		payloadRefJSON,
		metadataJSON,
		c.Version,
		ttlSeconds,
		accessControlJSON,
		expiresAt,
		c.Checksum,
	)

	if result.RowsAffected() == 0 {
		return fmt.Errorf("context not found")
	}
	return err
}

// Delete removes a context
func (s *PostgresContextStore) Delete(ctx context.Context, contextID string) error {
	query := `DELETE FROM contexts WHERE id = $1`
	result, err := s.pool.Exec(ctx, query, contextID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("context not found")
	}
	return nil
}

// List retrieves contexts with filters
func (s *PostgresContextStore) List(ctx context.Context, filters *storage.ContextFilters) ([]*models.Context, error) {
	query := `SELECT id, type, agent_id, tenant_id, payload, payload_ref, metadata, version, schema_id, ttl_seconds, access_control, created_at, expires_at, checksum FROM contexts WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filters == nil {
		filters = &storage.ContextFilters{Limit: 100}
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

	if filters.AgentID != "" {
		query += fmt.Sprintf(" AND agent_id = $%d", argIndex)
		args = append(args, filters.AgentID)
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

	var contexts []*models.Context
	for rows.Next() {
		var c models.Context
		var payloadRefJSON []byte
		var metadataJSON []byte
		var accessControlJSON []byte
		var ttlSeconds int
		var expiresAt sql.NullTime

		err := rows.Scan(
			&c.ID,
			&c.Type,
			&c.AgentID,
			&c.TenantID,
			&c.Payload,
			&payloadRefJSON,
			&metadataJSON,
			&c.Version,
			&c.SchemaID,
			&ttlSeconds,
			&accessControlJSON,
			&c.CreatedAt,
			&expiresAt,
			&c.Checksum,
		)
		if err != nil {
			return nil, err
		}

		c.TTL = time.Duration(ttlSeconds) * time.Second
		if expiresAt.Valid {
			c.ExpiresAt = expiresAt.Time
		}

		if err := json.Unmarshal(metadataJSON, &c.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
		if err := json.Unmarshal(accessControlJSON, &c.AccessControl); err != nil {
			return nil, fmt.Errorf("failed to unmarshal access_control: %w", err)
		}
		if len(payloadRefJSON) > 0 {
			if err := json.Unmarshal(payloadRefJSON, &c.PayloadRef); err != nil {
				return nil, fmt.Errorf("failed to unmarshal payload_ref: %w", err)
			}
		}

		contexts = append(contexts, &c)
	}

	return contexts, rows.Err()
}

// DeleteExpired removes expired contexts
func (s *PostgresContextStore) DeleteExpired(ctx context.Context) (int, error) {
	query := `DELETE FROM contexts WHERE expires_at IS NOT NULL AND expires_at < NOW()`
	result, err := s.pool.Exec(ctx, query)
	if err != nil {
		return 0, err
	}
	return int(result.RowsAffected()), nil
}

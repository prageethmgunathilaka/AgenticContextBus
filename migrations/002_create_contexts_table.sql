-- Migration: 002_create_contexts_table.sql
-- Creates the contexts table for context management

CREATE TABLE IF NOT EXISTS contexts (
    id VARCHAR(255) PRIMARY KEY,
    type VARCHAR(100) NOT NULL,
    agent_id VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255) NOT NULL DEFAULT 'default', -- Single tenant in MVP
    payload BYTEA, -- For small contexts <1MB
    payload_ref JSONB, -- For large contexts (Phase 2+)
    metadata JSONB,
    version VARCHAR(50) NOT NULL DEFAULT '1.0',
    schema_id VARCHAR(255),
    ttl_seconds INTEGER DEFAULT 86400, -- 24 hours default
    access_control JSONB NOT NULL, -- {scope, allowed_ids}
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    checksum VARCHAR(64), -- SHA-256
    
    CONSTRAINT contexts_type_check CHECK (type IS NOT NULL AND type != ''),
    CONSTRAINT contexts_access_control_check CHECK (access_control IS NOT NULL)
);

CREATE INDEX IF NOT EXISTS idx_contexts_tenant_id ON contexts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_contexts_agent_id ON contexts(agent_id);
CREATE INDEX IF NOT EXISTS idx_contexts_type ON contexts(type);
CREATE INDEX IF NOT EXISTS idx_contexts_expires_at ON contexts(expires_at) WHERE expires_at IS NOT NULL;


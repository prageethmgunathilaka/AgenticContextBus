-- Migration: 001_create_agents_table.sql
-- Creates the agents table for agent registry

CREATE TABLE IF NOT EXISTS agents (
    id VARCHAR(255) PRIMARY KEY,
    type VARCHAR(100) NOT NULL,
    location VARCHAR(255),
    capabilities TEXT[], -- Array of strings
    metadata JSONB,
    status VARCHAR(50) NOT NULL DEFAULT 'unknown',
    tenant_id VARCHAR(255) NOT NULL DEFAULT 'default', -- Single tenant in MVP
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_seen TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT agents_status_check CHECK (status IN ('online', 'offline', 'unknown'))
);

CREATE INDEX IF NOT EXISTS idx_agents_tenant_id ON agents(tenant_id);
CREATE INDEX IF NOT EXISTS idx_agents_type ON agents(type);
CREATE INDEX IF NOT EXISTS idx_agents_status ON agents(status);
CREATE INDEX IF NOT EXISTS idx_agents_last_seen ON agents(last_seen);


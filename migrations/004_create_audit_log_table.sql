-- Migration: 004_create_audit_log_table.sql
-- Creates the audit log table for security and compliance

CREATE TABLE IF NOT EXISTS audit_log (
    id SERIAL PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL DEFAULT 'default',
    agent_id VARCHAR(255),
    action VARCHAR(100) NOT NULL, -- register, unregister, login, context_read, context_write, etc.
    resource_type VARCHAR(50), -- agent, context, message, etc.
    resource_id VARCHAR(255),
    details JSONB,
    ip_address VARCHAR(45), -- IPv6 max length
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_log_tenant_id ON audit_log(tenant_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_agent_id ON audit_log(agent_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_action ON audit_log(action);
CREATE INDEX IF NOT EXISTS idx_audit_log_created_at ON audit_log(created_at);

-- Partition by month for better performance (optional, Phase 2+)
-- CREATE TABLE audit_log_2025_01 PARTITION OF audit_log FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');


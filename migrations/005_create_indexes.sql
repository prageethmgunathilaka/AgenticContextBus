-- Migration: 005_create_indexes.sql
-- Additional indexes for performance optimization

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_agents_tenant_status ON agents(tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_contexts_tenant_type ON contexts(tenant_id, type);
CREATE INDEX IF NOT EXISTS idx_contexts_agent_created ON contexts(agent_id, created_at DESC);

-- Index for message routing
CREATE INDEX IF NOT EXISTS idx_messages_topic_status ON messages(topic, status);

-- Index for audit log queries
CREATE INDEX IF NOT EXISTS idx_audit_log_tenant_action ON audit_log(tenant_id, action);


-- Migration: 003_create_messages_table.sql
-- Creates the messages table for message tracking and audit

CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(255) PRIMARY KEY,
    from_agent_id VARCHAR(255) NOT NULL,
    to_agent_id VARCHAR(255), -- NULL for broadcast
    topic VARCHAR(255) NOT NULL,
    context_id VARCHAR(255),
    message_type VARCHAR(50) NOT NULL,
    idempotency_key VARCHAR(255) NOT NULL,
    correlation_id VARCHAR(255),
    reply_to VARCHAR(255),
    metadata JSONB,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) NOT NULL DEFAULT 'queued', -- queued, delivered, failed
    
    CONSTRAINT messages_type_check CHECK (message_type IN ('event', 'command', 'query', 'response')),
    CONSTRAINT messages_status_check CHECK (status IN ('queued', 'delivered', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_messages_from_agent ON messages(from_agent_id);
CREATE INDEX IF NOT EXISTS idx_messages_to_agent ON messages(to_agent_id);
CREATE INDEX IF NOT EXISTS idx_messages_topic ON messages(topic);
CREATE INDEX IF NOT EXISTS idx_messages_idempotency_key ON messages(idempotency_key);
CREATE INDEX IF NOT EXISTS idx_messages_correlation_id ON messages(correlation_id);
CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp);


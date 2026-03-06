-- Events table in PostgreSQL for analytics
-- Requires: CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS events (
    event_id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       UUID        NOT NULL,
    user_id         TEXT        NOT NULL,
    session_id      TEXT        NOT NULL,
    event_name      TEXT        NOT NULL,
    event_timestamp TIMESTAMPTZ NOT NULL,
    received_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    properties      JSONB       NOT NULL DEFAULT '{}'::jsonb,
    context         JSONB       NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX IF NOT EXISTS idx_events_tenant_time
    ON events (tenant_id, event_timestamp);


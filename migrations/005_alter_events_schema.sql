-- Align events table with spec: id, nullable session_id, indexes
ALTER TABLE events RENAME COLUMN event_id TO id;

ALTER TABLE events ALTER COLUMN session_id DROP NOT NULL;

CREATE INDEX IF NOT EXISTS idx_events_tenant_user
    ON events (tenant_id, user_id);

CREATE INDEX IF NOT EXISTS idx_events_tenant_event_name
    ON events (tenant_id, event_name);

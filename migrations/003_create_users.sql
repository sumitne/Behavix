-- Users: end users of the SaaS product (lifecycle analytics)
CREATE TABLE IF NOT EXISTS users (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         UUID NOT NULL,
    external_user_id  TEXT NOT NULL,
    email             TEXT,
    created_at        TIMESTAMPTZ DEFAULT NOW(),
    metadata          JSONB
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_tenant_external
    ON users (tenant_id, external_user_id);

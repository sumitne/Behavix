-- Insights: generated behavioral insights
CREATE TABLE IF NOT EXISTS insights (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      UUID NOT NULL,
    type           TEXT NOT NULL,
    title          TEXT NOT NULL,
    description    TEXT,
    severity       TEXT,
    metric_value   FLOAT,
    baseline_value FLOAT,
    metadata       JSONB,
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_insights_tenant_created
    ON insights (tenant_id, created_at DESC);

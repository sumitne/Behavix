-- Seed data for e2e: one tenant, sample events, sample insights.
-- API key: dev-api-key-12345

INSERT INTO tenants (name, api_key)
VALUES ('Acme Corp (Dev)', 'dev-api-key-12345')
ON CONFLICT (api_key) DO UPDATE SET name = EXCLUDED.name;

INSERT INTO events (tenant_id, user_id, session_id, event_name, event_timestamp, received_at, properties, context)
SELECT t.id, 'user_001', 'sess_alpha', 'page_view', NOW() - INTERVAL '2 hours', NOW(), '{"path":"/dashboard"}'::jsonb, '{"platform":"web"}'::jsonb
FROM tenants t WHERE t.api_key = 'dev-api-key-12345';

INSERT INTO events (tenant_id, user_id, session_id, event_name, event_timestamp, received_at, properties, context)
SELECT t.id, 'user_001', 'sess_alpha', 'project_created', NOW() - INTERVAL '90 minutes', NOW(), '{"project_type":"demo"}'::jsonb, '{"platform":"web"}'::jsonb
FROM tenants t WHERE t.api_key = 'dev-api-key-12345';

INSERT INTO events (tenant_id, user_id, session_id, event_name, event_timestamp, received_at, properties, context)
SELECT t.id, 'user_002', 'sess_beta', 'signup', NOW() - INTERVAL '1 hour', NOW(), '{"source":"organic"}'::jsonb, '{"platform":"mobile"}'::jsonb
FROM tenants t WHERE t.api_key = 'dev-api-key-12345';

INSERT INTO events (tenant_id, user_id, session_id, event_name, event_timestamp, received_at, properties, context)
SELECT t.id, 'user_002', NULL, 'feature_used', NOW() - INTERVAL '30 minutes', NOW(), '{"feature":"export"}'::jsonb, '{}'::jsonb
FROM tenants t WHERE t.api_key = 'dev-api-key-12345';

INSERT INTO events (tenant_id, user_id, session_id, event_name, event_timestamp, received_at, properties, context)
SELECT t.id, 'user_003', 'sess_gamma', 'page_view', NOW() - INTERVAL '15 minutes', NOW(), '{"path":"/settings"}'::jsonb, '{"platform":"web"}'::jsonb
FROM tenants t WHERE t.api_key = 'dev-api-key-12345';

INSERT INTO insights (tenant_id, type, title, description, severity, metric_value, baseline_value, metadata)
SELECT t.id, 'activation_drop', 'Activation dropped 31%', 'Activation dropped significantly compared to the 7-day average.', 'high', 0.29, 0.42, '{}'::jsonb
FROM tenants t WHERE t.api_key = 'dev-api-key-12345';

INSERT INTO insights (tenant_id, type, title, description, severity, metric_value, baseline_value, metadata)
SELECT t.id, 'feature_adoption', 'Export feature usage up 15%', 'More users are using the export feature this week.', 'medium', 0.45, 0.39, '{}'::jsonb
FROM tenants t WHERE t.api_key = 'dev-api-key-12345';

INSERT INTO insights (tenant_id, type, title, description, severity, metric_value, baseline_value, metadata)
SELECT t.id, 'retention_slip', 'Day-7 retention slightly down', 'Week-over-week day-7 retention decreased by 4%.', 'low', 0.62, 0.65, '{}'::jsonb
FROM tenants t WHERE t.api_key = 'dev-api-key-12345';

package postgres

import (
	"context"
	"encoding/json"

	"behavix-ai/internal/domain/insight"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InsightRepository implements insight.Repository using PostgreSQL.
type InsightRepository struct {
	pool *pgxpool.Pool
}

// NewInsightRepository returns an insight repository backed by PostgreSQL.
func NewInsightRepository(pool *pgxpool.Pool) *InsightRepository {
	return &InsightRepository{pool: pool}
}

// List returns recent insights for a tenant.
func (r *InsightRepository) List(ctx context.Context, tenantID uuid.UUID, limit int) ([]insight.Insight, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	rows, err := r.pool.Query(ctx,
		`SELECT id, tenant_id, type, title, description, severity, metric_value, baseline_value, metadata, created_at
		 FROM insights WHERE tenant_id = $1 ORDER BY created_at DESC LIMIT $2`,
		tenantID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []insight.Insight
	for rows.Next() {
		var in insight.Insight
		var desc, sev *string
		var meta []byte
		var metricVal, baselineVal *float64
		err := rows.Scan(&in.ID, &in.TenantID, &in.Type, &in.Title, &desc, &sev, &metricVal, &baselineVal, &meta, &in.CreatedAt)
		if err != nil {
			return nil, err
		}
		if desc != nil {
			in.Description = *desc
		}
		if sev != nil {
			in.Severity = *sev
		}
		in.MetricValue = metricVal
		in.BaselineValue = baselineVal
		_ = json.Unmarshal(meta, &in.Metadata)
		if in.Metadata == nil {
			in.Metadata = make(map[string]interface{})
		}
		out = append(out, in)
	}
	return out, rows.Err()
}

// Compile-time check that InsightRepository implements insight.Repository.
var _ insight.Repository = (*InsightRepository)(nil)

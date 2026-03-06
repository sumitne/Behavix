package postgres

import (
	"context"
	"errors"

	"behavix-ai/internal/domain/tenant"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrTenantNotFound is returned when no tenant matches the API key.
var ErrTenantNotFound = errors.New("tenant not found")

// TenantRepository implements tenant.Repository using PostgreSQL.
type TenantRepository struct {
	pool *pgxpool.Pool
}

// NewTenantRepository returns a tenant repository backed by PostgreSQL.
func NewTenantRepository(pool *pgxpool.Pool) *TenantRepository {
	return &TenantRepository{pool: pool}
}

// GetTenantIDByAPIKey returns the tenant ID for the given API key.
func (r *TenantRepository) GetTenantIDByAPIKey(ctx context.Context, apiKey string) (uuid.UUID, error) {
	t, err := r.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return uuid.Nil, err
	}
	return t.ID, nil
}

// GetByAPIKey returns the tenant for the given API key.
func (r *TenantRepository) GetByAPIKey(ctx context.Context, apiKey string) (*tenant.Tenant, error) {
	var t tenant.Tenant
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, api_key, created_at FROM tenants WHERE api_key = $1`,
		apiKey,
	).Scan(&t.ID, &t.Name, &t.APIKey, &t.CreatedAt)
	if err != nil {
		return nil, ErrTenantNotFound
	}
	return &t, nil
}

// Compile-time check that TenantRepository implements tenant.Repository.
var _ tenant.Repository = (*TenantRepository)(nil)

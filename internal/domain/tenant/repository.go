package tenant

import (
	"context"

	"github.com/google/uuid"
)

// Repository provides tenant data access (port in clean architecture).
type Repository interface {
	GetTenantIDByAPIKey(ctx context.Context, apiKey string) (uuid.UUID, error)
	GetByAPIKey(ctx context.Context, apiKey string) (*Tenant, error)
}

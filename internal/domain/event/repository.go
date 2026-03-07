package event

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TenantStats holds aggregate stats for a tenant's events.
type TenantStats struct {
	ActiveUsers  int64
	EventsCount  int64
	LastActivity *time.Time
}

// Repository persists and queries events (port in clean architecture).
type Repository interface {
	InsertBatch(ctx context.Context, records []Record) error
	List(ctx context.Context, tenantID uuid.UUID, limit int, eventName, userID string) ([]Record, error)
	TenantStats(ctx context.Context, tenantID uuid.UUID) (TenantStats, error)
}

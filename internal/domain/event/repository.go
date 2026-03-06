package event

import (
	"context"

	"github.com/google/uuid"
)

// Repository persists and queries events (port in clean architecture).
type Repository interface {
	InsertBatch(ctx context.Context, records []Record) error
	List(ctx context.Context, tenantID uuid.UUID, limit int, eventName, userID string) ([]Record, error)
}

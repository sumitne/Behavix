package insight

import (
	"context"

	"github.com/google/uuid"
)

// Repository provides insight data access (port in clean architecture).
type Repository interface {
	List(ctx context.Context, tenantID uuid.UUID, limit int) ([]Insight, error)
}

package insight

import (
	"context"

	"behavix-ai/internal/domain/insight"

	"github.com/google/uuid"
)

// Service provides insight use cases (e.g. feed for dashboard).
type Service struct {
	insightRepo insight.Repository
}

// NewService creates an insight service.
func NewService(insightRepo insight.Repository) *Service {
	return &Service{insightRepo: insightRepo}
}

// List returns recent insights for the tenant.
func (s *Service) List(ctx context.Context, tenantID uuid.UUID, limit int) ([]insight.Insight, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.insightRepo.List(ctx, tenantID, limit)
}

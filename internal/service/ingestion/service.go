package ingestion

import (
	"context"
	"errors"
	"time"

	"behavix-ai/internal/domain/event"

	"github.com/google/uuid"
)

var (
	ErrEmptyEvents   = errors.New("events list is empty")
	ErrInvalidEvent  = errors.New("invalid event: user_id and event_name required")
)

// Service handles event ingestion use cases: validate, enrich, persist.
type Service struct {
	eventRepo event.Repository
}

// NewService creates an event ingestion service.
func NewService(eventRepo event.Repository) *Service {
	return &Service{eventRepo: eventRepo}
}

// IngestBatch validates events, attaches tenant_id, generates id and received_at, and inserts.
func (s *Service) IngestBatch(ctx context.Context, tenantID uuid.UUID, events []event.Event) error {
	if len(events) == 0 {
		return ErrEmptyEvents
	}
	records := make([]event.Record, 0, len(events))
	now := time.Now().UTC()
	for i := range events {
		e := &events[i]
		if e.UserID == "" || e.EventName == "" {
			return ErrInvalidEvent
		}
		rec := event.Record{
			ID:            uuid.New(),
			TenantID:      tenantID,
			UserID:        e.UserID,
			SessionID:     e.SessionID,
			EventName:     e.EventName,
			EventTimestamp: e.EventTimestamp,
			ReceivedAt:    now,
			Properties:    e.Properties,
			Context:       e.Context,
		}
		if rec.Properties == nil {
			rec.Properties = make(map[string]interface{})
		}
		if rec.Context == nil {
			rec.Context = make(map[string]interface{})
		}
		records = append(records, rec)
	}
	return s.eventRepo.InsertBatch(ctx, records)
}

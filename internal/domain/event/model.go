package event

import (
	"time"

	"github.com/google/uuid"
)

// Event represents a single product usage event (API input).
type Event struct {
	UserID        string                 `json:"user_id"`
	SessionID     string                 `json:"session_id"`
	EventName     string                 `json:"event_name"`
	EventTimestamp time.Time             `json:"event_timestamp"`
	Properties    map[string]interface{} `json:"properties"`
	Context       map[string]interface{} `json:"context"`
}

// BatchRequest is the payload for POST /api/v1/events.
type BatchRequest struct {
	Events []Event `json:"events"`
}

// Record is a persisted event (with id, tenant_id, received_at).
type Record struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	UserID        string
	SessionID     string
	EventName     string
	EventTimestamp time.Time
	ReceivedAt    time.Time
	Properties    map[string]interface{}
	Context       map[string]interface{}
}

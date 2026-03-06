package insight

import (
	"time"

	"github.com/google/uuid"
)

// Insight represents a generated behavioral insight.
type Insight struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	Type          string
	Title         string
	Description   string
	Severity      string
	MetricValue   *float64
	BaselineValue *float64
	Metadata      map[string]interface{}
	CreatedAt     time.Time
}

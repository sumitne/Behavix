package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents an end user of the SaaS product (lifecycle analytics).
type User struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	ExternalUserID   string
	Email            string
	CreatedAt        time.Time
	Metadata         map[string]interface{}
}

package tenant

import (
	"time"

	"github.com/google/uuid"
)

// Tenant represents a tenant in the system.
type Tenant struct {
	ID        uuid.UUID
	Name      string
	APIKey    string
	CreatedAt time.Time
}

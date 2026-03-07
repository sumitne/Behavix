package customer

import (
	"net/http"

	"behavix-ai/internal/domain/event"
	"behavix-ai/internal/domain/tenant"
	httpserver "behavix-ai/internal/infrastructure/http/server"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles GET /api/v1/customers and GET /api/v1/customers/:id.
type Handler struct {
	tenantRepo tenant.Repository
	eventRepo  event.Repository
}

// NewHandler creates a customer HTTP handler.
func NewHandler(tenantRepo tenant.Repository, eventRepo event.Repository) *Handler {
	return &Handler{tenantRepo: tenantRepo, eventRepo: eventRepo}
}

// customerResponse matches the frontend Customer type.
type customerResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	HealthScore   int    `json:"health_score"`
	ActiveUsers   int64  `json:"active_users"`
	LastActivity  string `json:"last_activity"`
}

func buildCustomerResponse(t *tenant.Tenant, stats event.TenantStats) customerResponse {
	resp := customerResponse{
		ID:          t.ID.String(),
		Name:        t.Name,
		HealthScore: 0, // placeholder until we have health logic
		ActiveUsers: stats.ActiveUsers,
		LastActivity: "",
	}
	if stats.LastActivity != nil {
		resp.LastActivity = stats.LastActivity.Format("2006-01-02T15:04:05Z07:00")
	}
	return resp
}

// List returns the current tenant as the only customer (single-tenant view).
func (h *Handler) List(c *gin.Context) {
	tenantID, ok := httpserver.GetTenantIDFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}
	t, err := h.tenantRepo.GetByID(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load tenant"})
		return
	}
	stats, err := h.eventRepo.TenantStats(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load stats"})
		return
	}
	c.JSON(http.StatusOK, []customerResponse{buildCustomerResponse(t, stats)})
}

// Get returns a single customer by id if it matches the current tenant.
func (h *Handler) Get(c *gin.Context) {
	tenantID, ok := httpserver.GetTenantIDFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
		return
	}
	if id != tenantID {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	t, err := h.tenantRepo.GetByID(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load tenant"})
		return
	}
	stats, err := h.eventRepo.TenantStats(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load stats"})
		return
	}
	c.JSON(http.StatusOK, buildCustomerResponse(t, stats))
}

// usageResponse matches the frontend CustomerUsage type.
type usageResponse struct {
	CustomerID      string             `json:"customer_id"`
	Period          string             `json:"period"`
	ActiveUsers     int64              `json:"active_users"`
	EventsCount     int64              `json:"events_count"`
	FeatureAdoption map[string]float64 `json:"feature_adoption,omitempty"`
}

// Usage returns usage stats for the customer (tenant) by id.
func (h *Handler) Usage(c *gin.Context) {
	tenantID, ok := httpserver.GetTenantIDFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
		return
	}
	if id != tenantID {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	stats, err := h.eventRepo.TenantStats(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load usage"})
		return
	}
	resp := usageResponse{
		CustomerID:  id.String(),
		Period:      "30d",
		ActiveUsers: stats.ActiveUsers,
		EventsCount: stats.EventsCount,
	}
	c.JSON(http.StatusOK, resp)
}

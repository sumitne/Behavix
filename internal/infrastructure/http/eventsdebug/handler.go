package eventsdebug

import (
	"net/http"
	"strconv"

	"behavix-ai/internal/domain/event"
	httpserver "behavix-ai/internal/infrastructure/http/server"

	"github.com/gin-gonic/gin"
)

// Handler handles GET /api/v1/events for debugging (list events with optional filters).
type Handler struct {
	eventRepo event.Repository
}

// NewHandler creates an events debug HTTP handler.
func NewHandler(eventRepo event.Repository) *Handler {
	return &Handler{eventRepo: eventRepo}
}

// List returns events for the tenant with optional query params: limit, event_name, user_id.
func (h *Handler) List(c *gin.Context) {
	tenantID, ok := httpserver.GetTenantIDFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant not found"})
		return
	}
	limit := 20
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}
	eventName := c.Query("event_name")
	userID := c.Query("user_id")

	records, err := h.eventRepo.List(c.Request.Context(), tenantID, limit, eventName, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list events"})
		return
	}
	out := make([]gin.H, 0, len(records))
	for _, r := range records {
		out = append(out, gin.H{
			"id":             r.ID,
			"tenant_id":      r.TenantID,
			"user_id":        r.UserID,
			"session_id":     r.SessionID,
			"event_name":     r.EventName,
			"event_timestamp": r.EventTimestamp,
			"received_at":    r.ReceivedAt,
			"properties":     r.Properties,
			"context":        r.Context,
		})
	}
	c.JSON(http.StatusOK, out)
}

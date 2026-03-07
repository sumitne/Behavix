package insight

import (
	"net/http"
	"strconv"

	"behavix-ai/internal/domain/tenant"
	httpserver "behavix-ai/internal/infrastructure/http/server"
	"behavix-ai/internal/service/insight"

	"github.com/gin-gonic/gin"
)

// Handler handles GET /api/v1/insights.
type Handler struct {
	svc        *insight.Service
	tenantRepo tenant.Repository
}

// NewHandler creates an insight HTTP handler.
func NewHandler(svc *insight.Service, tenantRepo tenant.Repository) *Handler {
	return &Handler{svc: svc, tenantRepo: tenantRepo}
}

// List returns recent insights for the tenant.
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
	insights, err := h.svc.List(c.Request.Context(), tenantID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list insights"})
		return
	}
	// Resolve tenant name for customer_name in response
	customerName := ""
	if t, err := h.tenantRepo.GetByID(c.Request.Context(), tenantID); err == nil {
		customerName = t.Name
	}
	// Map to API response shape
	out := make([]gin.H, 0, len(insights))
	for _, in := range insights {
		item := gin.H{
			"id":            in.ID,
			"type":          in.Type,
			"title":         in.Title,
			"description":   in.Description,
			"severity":      in.Severity,
			"created_at":    in.CreatedAt,
			"customer_name": customerName,
		}
		if in.MetricValue != nil {
			item["metric_value"] = *in.MetricValue
		}
		if in.BaselineValue != nil {
			item["baseline_value"] = *in.BaselineValue
		}
		out = append(out, item)
	}
	c.JSON(http.StatusOK, out)
}

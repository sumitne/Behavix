package ingestion

import (
	"errors"
	"net/http"

	"behavix-ai/internal/domain/event"
	httpserver "behavix-ai/internal/infrastructure/http/server"
	"behavix-ai/internal/service/ingestion"
	"behavix-ai/pkg/response"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for event ingestion.
type Handler struct {
	svc *ingestion.Service
}

// NewHandler creates an ingestion HTTP handler.
func NewHandler(svc *ingestion.Service) *Handler {
	return &Handler{svc: svc}
}

// Ingest implements httpserver.EventIngestionHandler.
func (h *Handler) Ingest(c *gin.Context) {
	var req event.BatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid payload")
		return
	}
	tenantID, ok := httpserver.GetTenantIDFromContext(c.Request.Context())
	if !ok {
		response.Error(c, http.StatusUnauthorized, "tenant not found")
		return
	}
	if err := h.svc.IngestBatch(c.Request.Context(), tenantID, req.Events); err != nil {
		if errors.Is(err, ingestion.ErrEmptyEvents) || errors.Is(err, ingestion.ErrInvalidEvent) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, "ingestion failed")
		return
	}
	response.Accepted(c)
}

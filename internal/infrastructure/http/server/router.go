package server

import (
	"net/http"

	"behavix-ai/internal/domain/tenant"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RouterConfig holds dependencies for building the HTTP router.
type RouterConfig struct {
	Logger          *zap.Logger
	EventIngestion  EventIngestionHandler
	EventList       EventListHandler
	InsightList     InsightListHandler
	TenantRepo      tenant.Repository
	UseTenantAuth   bool
}

// EventIngestionHandler handles POST /api/v1/events.
type EventIngestionHandler interface {
	Ingest(c *gin.Context)
}

// EventListHandler handles GET /api/v1/events (debug).
type EventListHandler interface {
	List(c *gin.Context)
}

// InsightListHandler handles GET /api/v1/insights.
type InsightListHandler interface {
	List(c *gin.Context)
}

// NewRouter builds the Gin router with health, versioned API, and middleware.
func NewRouter(cfg RouterConfig) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(RequestLogger(cfg.Logger))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	if cfg.UseTenantAuth && cfg.TenantRepo != nil {
		v1.Use(TenantAuth(cfg.TenantRepo))
	}
	if cfg.EventIngestion != nil {
		v1.POST("/events", cfg.EventIngestion.Ingest)
	}
	if cfg.EventList != nil {
		v1.GET("/events", cfg.EventList.List)
	}
	if cfg.InsightList != nil {
		v1.GET("/insights", cfg.InsightList.List)
	}

	return r
}

package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"behavix-ai/internal/domain/tenant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contextKey string

const TenantIDKey contextKey = "tenant_id"

// RequestLogger logs HTTP requests with method, path, status, and latency.
func RequestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		c.Next()
		log.Info("request",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", clientIP),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}

// TenantAuth extracts the API key from Authorization: Bearer <API_KEY>,
// resolves the tenant via tenant.Repository, and attaches tenant_id to the request context.
func TenantAuth(repo tenant.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}
		apiKey := strings.TrimSpace(parts[1])
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing api key"})
			return
		}
		tenantID, err := repo.GetTenantIDByAPIKey(c.Request.Context(), apiKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			return
		}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), TenantIDKey, tenantID))
		c.Next()
	}
}

// GetTenantIDFromContext returns the tenant_id stored in the request context by TenantAuth.
func GetTenantIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	v, ok := ctx.Value(TenantIDKey).(uuid.UUID)
	return v, ok
}

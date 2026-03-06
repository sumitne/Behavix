package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Accepted sends a 202 response with status "accepted".
func Accepted(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
}

// Error sends a JSON error response with the given status and message.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

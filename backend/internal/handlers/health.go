package handlers

import (
	"net/http"

	"supernote-ai/backend/internal/database"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoint
type HealthHandler struct {
	db *database.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck returns the health status of the API
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Check database connection
	if err := h.db.HealthCheck(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "unhealthy",
			"database": "disconnected",
			"error":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "healthy",
		"database": "connected",
		"service":  "supernote-ai-backend",
		"version":  "1.0.0",
	})
}

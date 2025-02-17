package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Show the health status of the server.
// @Description Get the health status of the server.
// @Tags Health Status
// @Accept */*
// @Produce json
// @Success 200 {string} message "the server is healthy"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "the server is healthy",
	})
}

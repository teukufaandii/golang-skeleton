package handler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

func HealthCheck(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"uptime":    time.Since(startTime).String(),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
		"runtime": gin.H{
			"go_version":     runtime.Version(),
			"goroutines":     runtime.NumGoroutine(),
			"heap_alloc_mb":  float64(m.HeapAlloc) / 1024 / 1024,
			"total_alloc_mb": float64(m.TotalAlloc) / 1024 / 1024,
		},
	})
}

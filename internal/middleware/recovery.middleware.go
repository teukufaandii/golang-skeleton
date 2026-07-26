package middleware

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"golang-skeleton/pkg/utils"
)

// RecoveryMiddleware recovers from panics and logs the error
func RecoveryMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the stack trace
				logger.WithFields(logrus.Fields{
					"error":      err,
					"stacktrace": string(debug.Stack()),
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
				}).Error("Panic recovered")

				c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
					Success: false,
					Error:   "Internal server error",
					Code:    "INTERNAL_ERROR",
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}

// TimeoutMiddleware adds request timeout
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			c.JSON(http.StatusGatewayTimeout, utils.ErrorResponse{
				Success: false,
				Error:   "Request timed out",
				Code:    "GATEWAY_TIMEOUT",
			})
			c.Abort()
		}
	}
}

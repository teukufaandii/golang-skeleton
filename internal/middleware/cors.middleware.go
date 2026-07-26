package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns configured CORS middleware
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(config)
}

// DefaultCORSMiddleware for development - allows all origins
func DefaultCORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

package routes

import (
	handler "golang-skeleton/internal/handler/http"
	"golang-skeleton/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORSMiddleware([]string{"*"}))
	r.Use(middleware.LoggerMiddleware(logrus.New()))
	r.Use(middleware.RecoveryMiddleware(logrus.New()))
	r.Use(middleware.RateLimitMiddleware(20, 10))

	r.GET("/health", handler.HealthCheck)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)

			authProtected := auth.Group("")
			authProtected.Use(middleware.AuthMiddleware(jwtSecret))
			{
				// protected route
			}
		}
	}

	return r
}

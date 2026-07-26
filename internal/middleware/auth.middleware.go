package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"golang-skeleton/pkg/utils"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "user_id"
	UserEmailKey        = "user_email"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
				Success: false,
				Error:   "Authorization header is missing",
				Code:    "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
				Success: false,
				Error:   "Invalid Authorization Format",
				Code:    "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)

		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
				Success: false,
				Error:   err.Error(),
				Code:    "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Next()
	}
}

func AdminAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
				Success: false,
				Error:   "Unauthorized, header is missing",
				Code:    "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
				Success: false,
				Error:   "Invalid Authorization Format",
				Code:    "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerPrefix)

		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
				Success: false,
				Error:   err.Error(),
				Code:    "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		if claims.Role != "admin" {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse{
				Success: false,
				Error:   "Unauthorized, user is not admin",
				Code:    "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Next()
	}
}

func GetUserIDMiddleware(c *gin.Context) (uuid.UUID, bool) {
	userIDInterface, exists := c.Get(UserIDKey)
	if !exists {
		return uuid.UUID{}, false
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, false
	}

	return userID, true
}

func OptionalAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)

		if authHeader != "" && strings.HasPrefix(authHeader, BearerPrefix) {
			tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
			claims, err := utils.ValidateToken(tokenString, jwtSecret)
			if err == nil {
				c.Set(UserIDKey, claims.UserID)
				c.Set(UserEmailKey, claims.Email)
			}
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userIDInterface, exists := c.Get(UserIDKey)
	if !exists {
		return uuid.UUID{}, false
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {

		return uuid.UUID{}, false
	}

	return userID, true
}

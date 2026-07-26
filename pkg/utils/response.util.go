package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    Pagination  `json:"meta"`
}

// Pagination contains pagination metadata
type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// Success sends a successful response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 Created response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func BuildPaginatedResponse(data interface{}, page, limit int, total int64) *PaginatedResponse {
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: Pagination{
			CurrentPage: page,
			PerPage:     limit,
			TotalItems:  total,
			TotalPages:  totalPages,
		},
	}
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, data interface{}, page, limit int, total int64) {
	response := BuildPaginatedResponse(data, page, limit, total)
	c.JSON(http.StatusOK, response)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "BAD_REQUEST",
	})
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "UNAUTHORIZED",
	})
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "FORBIDDEN",
	})
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "NOT_FOUND",
	})
}

// Conflict sends a 409 Conflict response
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "CONFLICT",
	})
}

// TooManyRequests sends a 429 Too Many Requests response
func TooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "RATE_LIMIT_EXCEEDED",
	})
}

// InternalError sends a 500 Internal Server Error response
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "INTERNAL_ERROR",
	})
}

// ValidationError sends a 400 response with validation errors
func ValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Error:   err.Error(),
		Code:    "VALIDATION_ERROR",
	})
}

package middleware

import (
    "bytes"
    "io"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
)

const RequestIDKey = "X-Request-ID"

// LoggerMiddleware provides structured logging for all requests
func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Generate unique request ID
        requestID := c.GetHeader(RequestIDKey)
        if requestID == "" {
            requestID = uuid.New().String()
        }
        c.Set(RequestIDKey, requestID)
        c.Header(RequestIDKey, requestID)

        // Capture request details
        startTime := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery
        method := c.Request.Method
        clientIP := c.ClientIP()
        userAgent := c.Request.UserAgent()

        // Read body for logging (if not too large)
        var requestBody string
        if c.Request.Body != nil && c.Request.ContentLength < 10240 { // 10KB limit
            bodyBytes, _ := io.ReadAll(c.Request.Body)
            c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
            requestBody = string(bodyBytes)
        }

        // Process request
        c.Next()

        // Calculate latency
        latency := time.Since(startTime)
        statusCode := c.Writer.Status()

        // Log entry
        entry := logger.WithFields(logrus.Fields{
            "request_id": requestID,
            "status":     statusCode,
            "method":     method,
            "path":       path,
            "query":      query,
            "ip":         clientIP,
            "user_agent": userAgent,
            "latency_ms": latency.Milliseconds(),
            "latency":    latency.String(),
        })

        // Add user ID if authenticated
        if userID, exists := c.Get(UserIDKey); exists {
            entry = entry.WithField("user_id", userID)
        }

        // Log based on status code
        switch {
        case statusCode >= 500:
            entry.WithField("body", requestBody).Error("Server error")
        case statusCode >= 400:
            entry.WithField("body", requestBody).Warn("Client error")
        case statusCode >= 300:
            entry.Info("Redirect")
        default:
            entry.Info("Request completed")
        }
    }
}
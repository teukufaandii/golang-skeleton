package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"golang-skeleton/pkg/utils"
)

// IPRateLimiter stores rate limiters per IP
type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     *sync.RWMutex
	rate   rate.Limit
	burst  int
	ttl    time.Duration
	lastGC time.Time
}

// NewIPRateLimiter creates a new IP-based rate limiter
func NewIPRateLimiter(r rate.Limit, b int, ttl time.Duration) *IPRateLimiter {
	return &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		mu:     &sync.RWMutex{},
		rate:   r,
		burst:  b,
		ttl:    ttl,
		lastGC: time.Now(),
	}
}

// GetLimiter returns the rate limiter for the provided IP address
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Garbage collect old entries periodically
	if time.Since(i.lastGC) > i.ttl {
		i.ips = make(map[string]*rate.Limiter)
		i.lastGC = time.Now()
	}

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.rate, i.burst)
		i.ips[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware creates rate limiting middleware
// rate: requests per second, burst: max burst size
func RateLimitMiddleware(requestsPerSecond float64, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(rate.Limit(requestsPerSecond), burst, time.Hour)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.GetLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, utils.ErrorResponse{
				Success: false,
				Error:   "Too many requests",
				Code:    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// EndpointRateLimiter for specific endpoint rate limiting
type EndpointRateLimiter struct {
	limiters map[string]*IPRateLimiter
	mu       sync.RWMutex
}

func NewEndpointRateLimiter() *EndpointRateLimiter {
	return &EndpointRateLimiter{
		limiters: make(map[string]*IPRateLimiter),
	}
}

func (e *EndpointRateLimiter) ConfigureEndpoint(path string, rps float64, burst int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.limiters[path] = NewIPRateLimiter(rate.Limit(rps), burst, time.Hour)
}

func (e *EndpointRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		e.mu.RLock()
		limiter, exists := e.limiters[c.FullPath()]
		e.mu.RUnlock()

		if exists && !limiter.GetLimiter(c.ClientIP()).Allow() {
			c.JSON(http.StatusTooManyRequests, utils.ErrorResponse{
				Success: false,
				Error:   "Too many requests",
				Code:    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

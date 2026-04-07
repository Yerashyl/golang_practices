package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitInfo struct {
	Count      int
	LastAccess time.Time
}

var (
	rateLimitStore = make(map[string]*rateLimitInfo)
	mu             sync.Mutex
	Limit          = 5 // 5 requests per minute
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string

		// Try to identify by UserID from JWT first
		userID, exists := c.Get("userID")
		if exists {
			key = fmt.Sprintf("user:%v", userID)
		} else {
			key = fmt.Sprintf("ip:%s", c.ClientIP())
		}

		mu.Lock()
		info, ok := rateLimitStore[key]
		now := time.Now()

		if !ok || now.Sub(info.LastAccess) > time.Minute {
			// New window
			rateLimitStore[key] = &rateLimitInfo{
				Count:      1,
				LastAccess: now,
			}
			mu.Unlock()
			c.Next()
			return
		}

		if info.Count >= Limit {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests. Please try again later."})
			c.Abort()
			return
		}

		info.Count++
		mu.Unlock()
		c.Next()
	}
}

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
	Limit          = 10 // 10 requests per minute
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		var key string

		// If OptionalJWTAuthMiddleware set userID, use it. Otherwise use IP.
		userID, exists := c.Get("userID")
		if exists {
			key = fmt.Sprintf("user:%v", userID)
		} else {
			key = fmt.Sprintf("ip:%s", c.ClientIP())
		}

		mu.Lock()
		defer mu.Unlock()

		info, ok := rateLimitStore[key]
		now := time.Now()

		if !ok || now.Sub(info.LastAccess) > time.Minute {
			rateLimitStore[key] = &rateLimitInfo{
				Count:      1,
				LastAccess: now,
			}
			c.Next()
			return
		}

		if info.Count >= Limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Limit is " + fmt.Sprint(Limit) + " requests per minute.",
			})
			c.Abort()
			return
		}

		info.Count++
		c.Next()
	}
}

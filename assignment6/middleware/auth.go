package middleware

import (
	"assignment6/models"
	"assignment6/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func OptionalJWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := utils.ValidateToken(parts[1])
			if err == nil {
				c.Set("userID", claims.UserID)
				c.Set("email", claims.Email)
				c.Set("role", claims.Role)
			}
		}
		c.Next()
	}
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not found in context"})
			c.Abort()
			return
		}

		if role.(string) != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func VerifiedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		email, exists := c.Get("email")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User info not found"})
			c.Abort()
			return
		}

		models.Mu.RLock()
		user, ok := models.Users[email.(string)]
		models.Mu.RUnlock()

		if !ok || !user.IsVerified {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account not verified"})
			c.Abort()
			return
		}
		c.Next()
	}
}

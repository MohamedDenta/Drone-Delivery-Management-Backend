package middleware

import (
	"net/http"
	"strings"

	"github.com/MohamedDenta/Drone-Delivery-Management-Backend/internal/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store claims in context for handlers to use
		c.Set("user", claims.Name)
		c.Set("role", claims.UserType)

		c.Next()
	}
}

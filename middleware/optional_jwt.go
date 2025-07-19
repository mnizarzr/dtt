package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mnizarzr/dot-test/utils"
)

// OptionalJWTAuth middleware that allows both authenticated and unauthenticated requests
// If a valid JWT token is provided, it extracts user information
// If no token or invalid token, it continues without setting user context
func OptionalJWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, continue without user context
			c.Next()
			return
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			// Invalid format, continue without user context
			c.Next()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			// Empty token, continue without user context
			c.Next()
			return
		}

		// Validate token
		claims, err := utils.ValidateJWT(tokenString, jwtSecret)
		if err != nil {
			// Invalid token, continue without user context
			c.Next()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

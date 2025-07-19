package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mnizarzr/dot-test/common"
	"github.com/mnizarzr/dot-test/utils"
)

// JWTAuth middleware for JWT authentication
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.ErrorResponse(c, 401, "Authorization header required")
			c.Abort()
			return
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			common.ErrorResponse(c, 401, "Invalid authorization header format")
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			common.ErrorResponse(c, 401, "Token is required")
			c.Abort()
			return
		}

		// Validate token
		claims, err := utils.ValidateJWT(tokenString, jwtSecret)
		if err != nil {
			common.ErrorResponse(c, 401, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

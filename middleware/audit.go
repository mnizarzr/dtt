package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mnizarzr/dot-test/entity"
	"gorm.io/gorm"
)

// AuditResourceContext injects audit information into the request context for GORM hooks
func AuditResourceContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Inject user ID if available
		if userID, exists := c.Get("user_id"); exists {
			ctx = context.WithValue(ctx, entity.AuditUserIDKey, userID)
		}

		// Inject IP address
		ctx = context.WithValue(ctx, entity.AuditIPKey, c.ClientIP())

		// Inject user agent
		ctx = context.WithValue(ctx, entity.AuditUserAgent, c.Request.UserAgent())

		// Update request context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// AuditMiddleware creates audit logs for actions that don't modify database entities
// This is useful for actions like login, register, etc.
func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// after request processing, log the action if it was successful, something like bruteforce should be handled by other tools
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			path := c.Request.URL.Path
			method := c.Request.Method

			action := method + " " + path
			targetResource := method + " " + path

			ctx := c.Request.Context()

			if userID, exists := c.Get("user_id"); exists {
				ctx = context.WithValue(ctx, entity.AuditUserIDKey, userID)
			}

			ctx = context.WithValue(ctx, entity.AuditIPKey, c.ClientIP())
			ctx = context.WithValue(ctx, entity.AuditUserAgent, c.Request.UserAgent())
			targetID := uuid.Nil

			auditData := map[string]any{
				"method": method,
				"path":   path,
				"status": c.Writer.Status(),
			}

			entity.CreateAuditLog(db.WithContext(ctx), action, targetResource, targetID, auditData, nil)
		}
	}
}

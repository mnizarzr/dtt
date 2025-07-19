package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	auditLogTableName = "audit_logs"
)

type AuditLog struct {
	ID             uuid.UUID       `json:"id"`
	UserID         *uuid.UUID      `json:"user_id"`
	Action         string          `json:"action"`
	TargetResource string          `json:"target_resource"`
	TargetID       *uuid.UUID      `json:"target_id"`
	Details        json.RawMessage `json:"details"`
	CreatedAt      time.Time       `json:"created_at"`
}

func (*AuditLog) TableName() string {
	return auditLogTableName
}

// Audit action constants
const (
	AuditActionCreate = "CREATE"
	AuditActionUpdate = "UPDATE"
	AuditActionDelete = "DELETE"
)

// Context keys for audit information
type contextKey string

const (
	AuditUserIDKey contextKey = "audit_user_id"
	AuditIPKey     contextKey = "audit_ip"
	AuditUserAgent contextKey = "audit_user_agent"
)

// CreateAuditLog creates an audit log entry within the same transaction
func CreateAuditLog(tx *gorm.DB, action, resource string, targetID uuid.UUID, data interface{}, oldData interface{}) error {
	// Get audit context information
	var userID *uuid.UUID
	details := make(map[string]interface{})

	if ctx := tx.Statement.Context; ctx != nil {
		// Get user ID
		if uid := ctx.Value(AuditUserIDKey); uid != nil {
			if uidStr, ok := uid.(string); ok {
				if parsed, err := uuid.Parse(uidStr); err == nil {
					userID = &parsed
				}
			}
		}

		// Get IP address
		if ip := ctx.Value(AuditIPKey); ip != nil {
			if ipStr, ok := ip.(string); ok && ipStr != "" {
				details["ip_address"] = ipStr
			}
		}

		// Get user agent
		if ua := ctx.Value(AuditUserAgent); ua != nil {
			if uaStr, ok := ua.(string); ok && uaStr != "" {
				details["user_agent"] = uaStr
			}
		}
	}

	// Add data changes
	if data != nil {
		details["new"] = sanitizeForAudit(data)
	}
	if oldData != nil {
		details["old"] = sanitizeForAudit(oldData)
	}

	// Add timestamp
	details["timestamp"] = time.Now().UTC()

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		simpleDetails := map[string]interface{}{
			"action": action,
			"error":  "failed to serialize data",
		}
		detailsJSON, _ = json.Marshal(simpleDetails)
	}

	auditLog := &AuditLog{
		ID:             uuid.New(),
		UserID:         userID,
		Action:         action,
		TargetResource: resource,
		TargetID:       &targetID,
		Details:        detailsJSON,
		CreatedAt:      time.Now(),
	}

	return tx.Create(auditLog).Error
}

// sanitizeForAudit removes sensitive data before logging
func sanitizeForAudit(data interface{}) interface{} {
	// if data is user exclude PasswordHash
	if user, ok := data.(*User); ok {
		return map[string]interface{}{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		}
	}
	return data
}

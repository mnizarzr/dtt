package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	taskTableName = "tasks"
)

type Task struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	ProjectID   *uuid.UUID `json:"project_id"`
	CreatedBy   *uuid.UUID `json:"created_by"`
	AssignedTo  *uuid.UUID `json:"assigned_to"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (*Task) TableName() string {
	return taskTableName
}

// AfterCreate hook - logs task creation
func (t *Task) AfterCreate(tx *gorm.DB) error {
	return CreateAuditLog(tx, AuditActionCreate, taskTableName, t.ID, t, nil)
}

// AfterUpdate hook - logs task updates
func (t *Task) AfterUpdate(tx *gorm.DB) error {
	return CreateAuditLog(tx, AuditActionUpdate, taskTableName, t.ID, t, nil)
}

// AfterDelete hook - logs task deletion
func (t *Task) AfterDelete(tx *gorm.DB) error {
	return CreateAuditLog(tx, AuditActionDelete, taskTableName, t.ID, nil, t)
}

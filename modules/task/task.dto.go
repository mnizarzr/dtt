package task

import (
	"time"

	"github.com/google/uuid"
)

// Task status constants
const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
)

// Task priority constants
const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)

// CreateTaskRequest represents the request to create a new task
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=100"`
	Description string     `json:"description" binding:"max=1000"`
	Status      string     `json:"status" binding:"omitempty,oneof=pending in_progress completed"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	ProjectID   uuid.UUID  `json:"project_id" binding:"required"`
	AssignedTo  *uuid.UUID `json:"assigned_to,omitempty"`
}

// UpdateTaskRequest represents the request to update a task
type UpdateTaskRequest struct {
	Title       *string    `json:"title,omitempty" binding:"omitempty,min=1,max=100"`
	Description *string    `json:"description,omitempty" binding:"omitempty,max=1000"`
	Status      *string    `json:"status,omitempty" binding:"omitempty,oneof=pending in_progress completed"`
	Priority    *string    `json:"priority,omitempty" binding:"omitempty,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	AssignedTo  *uuid.UUID `json:"assigned_to,omitempty"`
}

// AssignTaskRequest represents the request to assign a task to a user
type AssignTaskRequest struct {
	AssignedTo uuid.UUID `json:"assigned_to" binding:"required"`
}

// TaskResponse represents a task in API responses
type TaskResponse struct {
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

// TaskListResponse represents a list of tasks with pagination
type TaskListResponse struct {
	Tasks []TaskResponse `json:"tasks"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

// TaskFilterRequest represents task filtering parameters
type TaskFilterRequest struct {
	ProjectID  *uuid.UUID `form:"project_id,omitempty"`
	AssignedTo *uuid.UUID `form:"assigned_to,omitempty"`
	Status     *string    `form:"status,omitempty" binding:"omitempty,oneof=pending in_progress completed"`
	Priority   *string    `form:"priority,omitempty" binding:"omitempty,oneof=low medium high"`
	Page       int        `form:"page" binding:"omitempty,min=1"`
	Limit      int        `form:"limit" binding:"omitempty,min=1,max=100"`
}

// SetDefaults sets default values for pagination
func (f *TaskFilterRequest) SetDefaults() {
	if f.Page == 0 {
		f.Page = 1
	}
	if f.Limit == 0 {
		f.Limit = 10
	}
}

// GetOffset calculates the offset for database queries
func (f *TaskFilterRequest) GetOffset() int {
	return (f.Page - 1) * f.Limit
}

// IsValidStatus checks if the status is valid
func IsValidStatus(status string) bool {
	return status == StatusPending || status == StatusInProgress || status == StatusCompleted
}

// IsValidPriority checks if the priority is valid
func IsValidPriority(priority string) bool {
	return priority == PriorityLow || priority == PriorityMedium || priority == PriorityHigh
}

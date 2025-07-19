package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/mnizarzr/dot-test/entity"
	"gorm.io/gorm"
)

// Repository defines the interface for task data operations
type Repository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID, offset, limit int) ([]*entity.Task, int64, error)
	GetByAssignedUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.Task, int64, error)
	GetWithFilters(ctx context.Context, filters TaskFilterRequest) ([]*entity.Task, int64, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetUserTasksInProject(ctx context.Context, userID, projectID uuid.UUID, offset, limit int) ([]*entity.Task, int64, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new task repository instance
func NewRepository(database *gorm.DB) Repository {
	return &repository{
		db: database,
	}
}

// Create creates a new task
func (r *repository) Create(ctx context.Context, task *entity.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// GetByID retrieves a task by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	var task entity.Task
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&task).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

// GetByProjectID retrieves tasks by project ID with pagination
func (r *repository) GetByProjectID(ctx context.Context, projectID uuid.UUID, offset, limit int) ([]*entity.Task, int64, error) {
	var tasks []*entity.Task
	var total int64

	// Count total tasks in project
	if err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("project_id = ?", projectID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get tasks with pagination
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, total, err
}

// GetByAssignedUser retrieves tasks assigned to a specific user with pagination
func (r *repository) GetByAssignedUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*entity.Task, int64, error) {
	var tasks []*entity.Task
	var total int64

	// Count total tasks assigned to user
	if err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("assigned_to = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get tasks with pagination
	err := r.db.WithContext(ctx).
		Where("assigned_to = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, total, err
}

// GetWithFilters retrieves tasks with various filters
func (r *repository) GetWithFilters(ctx context.Context, filters TaskFilterRequest) ([]*entity.Task, int64, error) {
	var tasks []*entity.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Task{})

	// Apply filters
	if filters.ProjectID != nil {
		query = query.Where("project_id = ?", *filters.ProjectID)
	}
	if filters.AssignedTo != nil {
		query = query.Where("assigned_to = ?", *filters.AssignedTo)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.Priority != nil {
		query = query.Where("priority = ?", *filters.Priority)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get tasks with pagination
	err := query.
		Offset(filters.GetOffset()).
		Limit(filters.Limit).
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, total, err
}

// Update updates an existing task
func (r *repository) Update(ctx context.Context, task *entity.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

// Delete deletes a task
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Task{}, id).Error
}

// GetUserTasksInProject gets tasks for a user within a specific project
func (r *repository) GetUserTasksInProject(ctx context.Context, userID, projectID uuid.UUID, offset, limit int) ([]*entity.Task, int64, error) {
	var tasks []*entity.Task
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Task{}).
		Where("assigned_to = ? AND project_id = ?", userID, projectID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get tasks with pagination
	err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&tasks).Error

	return tasks, total, err
}

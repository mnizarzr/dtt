package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/mnizarzr/dot-test/entity"
	"gorm.io/gorm"
)

// Repository defines the interface for project data operations
type Repository interface {
	Create(ctx context.Context, project *entity.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Project, error)
	GetAll(ctx context.Context, offset, limit int) ([]*entity.Project, int64, error)
	Update(ctx context.Context, project *entity.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByNameExcludingID(ctx context.Context, name string, excludeID uuid.UUID) (bool, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new project repository instance
func NewRepository(database *gorm.DB) Repository {
	return &repository{
		db: database,
	}
}

// Create creates a new project
func (r *repository) Create(ctx context.Context, project *entity.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

// GetByID retrieves a project by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Project, error) {
	var project entity.Project
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

// GetAll retrieves all projects with pagination
func (r *repository) GetAll(ctx context.Context, offset, limit int) ([]*entity.Project, int64, error) {
	var projects []*entity.Project
	var total int64

	// Count total projects
	if err := r.db.WithContext(ctx).Model(&entity.Project{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get projects with pagination
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&projects).Error

	return projects, total, err
}

// Update updates an existing project
func (r *repository) Update(ctx context.Context, project *entity.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

// Delete deletes a project (which will cascade delete tasks due to foreign key)
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Project{}, id).Error
}

// ExistsByName checks if a project with the given name exists
func (r *repository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Project{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// ExistsByNameExcludingID checks if a project with the given name exists, excluding a specific ID
func (r *repository) ExistsByNameExcludingID(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Project{}).
		Where("name = ? AND id != ?", name, excludeID).
		Count(&count).Error
	return count > 0, err
}

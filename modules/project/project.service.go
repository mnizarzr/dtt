package project

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mnizarzr/dot-test/common"
	"github.com/mnizarzr/dot-test/entity"
)

// Service defines the interface for project business logic
type Service interface {
	CreateProject(ctx context.Context, req CreateProjectRequest, createdBy uuid.UUID) (*ProjectResponse, error)
	GetProject(ctx context.Context, id uuid.UUID) (*ProjectResponse, error)
	GetAllProjects(ctx context.Context, pagination PaginationRequest) (*ProjectListResponse, error)
	UpdateProject(ctx context.Context, id uuid.UUID, req UpdateProjectRequest, userID uuid.UUID, userRole string) (*ProjectResponse, error)
	DeleteProject(ctx context.Context, id uuid.UUID, userID uuid.UUID, userRole string) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new project service instance
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// CreateProject creates a new project (only managers and admins can create projects)
func (s *service) CreateProject(ctx context.Context, req CreateProjectRequest, createdBy uuid.UUID) (*ProjectResponse, error) {
	// Check if project with same name already exists
	exists, err := s.repo.ExistsByName(ctx, req.Name)
	if err != nil {
		return nil, common.ErrFailedToCheckProject
	}
	if exists {
		return nil, common.ErrProjectNameAlreadyExists
	}

	// Create project entity
	now := time.Now()
	project := &entity.Project{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Save to database
	if err := s.repo.Create(ctx, project); err != nil {
		return nil, common.ErrFailedToCreateProject
	}

	return s.entityToResponse(project), nil
}

// GetProject retrieves a project by ID
func (s *service) GetProject(ctx context.Context, id uuid.UUID) (*ProjectResponse, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, common.ErrFailedToRetrieveProject
	}
	if project == nil {
		return nil, common.ErrProjectNotFound
	}

	return s.entityToResponse(project), nil
}

// GetAllProjects retrieves all projects with pagination
func (s *service) GetAllProjects(ctx context.Context, pagination PaginationRequest) (*ProjectListResponse, error) {
	pagination.SetDefaults()

	projects, total, err := s.repo.GetAll(ctx, pagination.GetOffset(), pagination.Limit)
	if err != nil {
		return nil, common.ErrFailedToRetrieveProjects
	}

	// Convert to response DTOs
	projectResponses := make([]ProjectResponse, len(projects))
	for i, project := range projects {
		projectResponses[i] = *s.entityToResponse(project)
	}

	return &ProjectListResponse{
		Projects: projectResponses,
		Total:    total,
		Page:     pagination.Page,
		Limit:    pagination.Limit,
	}, nil
}

// UpdateProject updates an existing project (only managers, admins, or creator can update)
func (s *service) UpdateProject(ctx context.Context, id uuid.UUID, req UpdateProjectRequest, userID uuid.UUID, userRole string) (*ProjectResponse, error) {
	// Get existing project
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, common.ErrFailedToRetrieveProject
	}
	if project == nil {
		return nil, common.ErrProjectNotFound
	}

	// Check permissions: only admin, manager, or creator can update
	if !s.canManageProject(userRole, userID, project.CreatedBy) {
		return nil, common.ErrForbidden
	}

	// Update fields if provided
	if req.Name != nil {
		// Check if new name conflicts with existing projects
		exists, err := s.repo.ExistsByNameExcludingID(ctx, *req.Name, id)
		if err != nil {
			return nil, common.ErrFailedToCheckProject
		}
		if exists {
			return nil, common.ErrProjectNameAlreadyExists
		}
		project.Name = *req.Name
	}

	if req.Description != nil {
		project.Description = *req.Description
	}

	project.UpdatedAt = time.Now()

	// Save changes
	if err := s.repo.Update(ctx, project); err != nil {
		return nil, common.ErrFailedToUpdateProject
	}

	return s.entityToResponse(project), nil
}

// DeleteProject deletes a project (only managers, admins, or creator can delete)
func (s *service) DeleteProject(ctx context.Context, id uuid.UUID, userID uuid.UUID, userRole string) error {
	// Get existing project
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return common.ErrFailedToRetrieveProject
	}
	if project == nil {
		return common.ErrProjectNotFound
	}

	// Check permissions: only admin, manager, or creator can delete
	if !s.canManageProject(userRole, userID, project.CreatedBy) {
		return common.ErrForbidden
	}

	// Delete project (tasks will be cascade deleted)
	if err := s.repo.Delete(ctx, id); err != nil {
		return common.ErrFailedToDeleteProject
	}

	return nil
}

// canManageProject checks if user can manage (update/delete) a project
func (s *service) canManageProject(userRole string, userID uuid.UUID, createdBy *uuid.UUID) bool {
	// Admin can manage any project
	if userRole == "admin" {
		return true
	}

	// Manager can manage any project
	if userRole == "manager" {
		return true
	}

	// Creator can manage their own project
	if createdBy != nil && *createdBy == userID {
		return true
	}

	return false
}

// entityToResponse converts a project entity to response DTO
func (s *service) entityToResponse(project *entity.Project) *ProjectResponse {
	return &ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		CreatedBy:   project.CreatedBy,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}
}

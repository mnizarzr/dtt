package task

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mnizarzr/dot-test/common"
	"github.com/mnizarzr/dot-test/entity"
	"github.com/mnizarzr/dot-test/modules/project"
	"github.com/mnizarzr/dot-test/modules/user"
	"github.com/mnizarzr/dot-test/utils"
)

// Service defines the interface for task business logic
type Service interface {
	CreateTask(ctx context.Context, req CreateTaskRequest, createdBy uuid.UUID, userRole string) (*TaskResponse, error)
	GetTask(ctx context.Context, id uuid.UUID, userID uuid.UUID, userRole string) (*TaskResponse, error)
	GetTasksWithFilters(ctx context.Context, filters TaskFilterRequest, userID uuid.UUID, userRole string) (*TaskListResponse, error)
	UpdateTask(ctx context.Context, id uuid.UUID, req UpdateTaskRequest, userID uuid.UUID, userRole string) (*TaskResponse, error)
	AssignTask(ctx context.Context, taskID uuid.UUID, req AssignTaskRequest, userID uuid.UUID, userRole string) (*TaskResponse, error)
	DeleteTask(ctx context.Context, id uuid.UUID, userID uuid.UUID, userRole string) error
}

// service implements the Service interface
type service struct {
	repo           Repository
	projectService project.Service
	userService    user.Service
}

// NewService creates a new task service instance
func NewService(repo Repository, projectService project.Service, userService user.Service) Service {
	return &service{
		repo:           repo,
		projectService: projectService,
		userService:    userService,
	}
}

// CreateTask creates a new task (must belong to a project)
func (s *service) CreateTask(ctx context.Context, req CreateTaskRequest, createdBy uuid.UUID, userRole string) (*TaskResponse, error) {
	// Verify project exists
	_, err := s.projectService.GetProject(ctx, req.ProjectID)
	if err != nil {
		if err == common.ErrProjectNotFound {
			return nil, common.ErrProjectNotFound
		}
		return nil, common.ErrFailedToRetrieveProject
	}

	// If assigning to someone, verify user exists and check permissions
	if req.AssignedTo != nil {
		// Only managers and admins can assign tasks
		if userRole != "admin" && userRole != "manager" {
			return nil, common.ErrCannotAssignTask
		}

		// Verify assigned user exists
		_, err := s.userService.GetUserByID(ctx, *req.AssignedTo)
		if err != nil {
			return nil, common.ErrUserNotFound
		}
	}

	// Set defaults
	status := req.Status
	if status == "" {
		status = StatusPending
	}

	priority := req.Priority
	if priority == "" {
		priority = PriorityMedium
	}

	// Adjust due date if it falls on a holiday
	dueDate := req.DueDate
	if dueDate != nil {
		adjustedDate := s.adjustForHolidays(*dueDate)
		dueDate = &adjustedDate
	}

	// Create task entity
	now := time.Now()
	task := &entity.Task{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		Priority:    priority,
		DueDate:     dueDate,
		ProjectID:   &req.ProjectID,
		CreatedBy:   &createdBy,
		AssignedTo:  req.AssignedTo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Save to database
	if err := s.repo.Create(ctx, task); err != nil {
		return nil, common.ErrFailedToCreateTask
	}

	return s.entityToResponse(task), nil
}

// GetTask retrieves a task by ID with permission checks
func (s *service) GetTask(ctx context.Context, id uuid.UUID, userID uuid.UUID, userRole string) (*TaskResponse, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, common.ErrFailedToRetrieveTask
	}
	if task == nil {
		return nil, common.ErrTaskNotFound
	}

	// Check permissions
	if !s.canViewTask(task, userID, userRole) {
		return nil, common.ErrForbidden
	}

	return s.entityToResponse(task), nil
}

// GetTasksWithFilters retrieves tasks with filters and permission checks
func (s *service) GetTasksWithFilters(ctx context.Context, filters TaskFilterRequest, userID uuid.UUID, userRole string) (*TaskListResponse, error) {
	filters.SetDefaults()

	switch userRole {
	case "admin":
	case "manager":
	case "user":
		filters.AssignedTo = &userID
	default:
		filters.AssignedTo = &userID
	}

	tasks, total, err := s.repo.GetWithFilters(ctx, filters)
	if err != nil {
		return nil, common.ErrFailedToRetrieveTasks
	}

	taskResponses := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = *s.entityToResponse(task)
	}

	return &TaskListResponse{
		Tasks: taskResponses,
		Total: total,
		Page:  filters.Page,
		Limit: filters.Limit,
	}, nil
}

// UpdateTask updates an existing task with permission checks
func (s *service) UpdateTask(ctx context.Context, id uuid.UUID, req UpdateTaskRequest, userID uuid.UUID, userRole string) (*TaskResponse, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, common.ErrFailedToRetrieveTask
	}
	if task == nil {
		return nil, common.ErrTaskNotFound
	}

	if !s.canUpdateTask(task, userID, userRole) {
		return nil, common.ErrForbidden
	}

	if userRole == "user" && task.AssignedTo != nil && *task.AssignedTo == userID {
		if req.Status != nil {
			task.Status = *req.Status
		}
	} else {
		if req.Title != nil {
			task.Title = *req.Title
		}
		if req.Description != nil {
			task.Description = *req.Description
		}
		if req.Status != nil {
			task.Status = *req.Status
		}
		if req.Priority != nil {
			task.Priority = *req.Priority
		}
		if req.DueDate != nil {
			task.DueDate = req.DueDate
		}
		if req.AssignedTo != nil {
			if userRole != "admin" && userRole != "manager" {
				return nil, common.ErrCannotAssignTask
			}
			_, err := s.userService.GetUserByID(ctx, *req.AssignedTo)
			if err != nil {
				return nil, common.ErrUserNotFound
			}
			task.AssignedTo = req.AssignedTo
		}
	}

	task.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, common.ErrFailedToUpdateTask
	}

	return s.entityToResponse(task), nil
}

// AssignTask assigns a task to a user (managers and admins only)
func (s *service) AssignTask(ctx context.Context, taskID uuid.UUID, req AssignTaskRequest, userID uuid.UUID, userRole string) (*TaskResponse, error) {
	if userRole != "admin" && userRole != "manager" {
		return nil, common.ErrCannotAssignTask
	}

	task, err := s.repo.GetByID(ctx, taskID)
	if err != nil {
		return nil, common.ErrFailedToRetrieveTask
	}
	if task == nil {
		return nil, common.ErrTaskNotFound
	}

	_, err = s.userService.GetUserByID(ctx, req.AssignedTo)
	if err != nil {
		return nil, common.ErrUserNotFound
	}

	task.AssignedTo = &req.AssignedTo
	task.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, common.ErrFailedToUpdateTask
	}

	return s.entityToResponse(task), nil
}

// DeleteTask deletes a task (only creator, manager, or admin can delete)
func (s *service) DeleteTask(ctx context.Context, id uuid.UUID, userID uuid.UUID, userRole string) error {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return common.ErrFailedToRetrieveTask
	}
	if task == nil {
		return common.ErrTaskNotFound
	}

	if !s.canDeleteTask(task, userID, userRole) {
		return common.ErrForbidden
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return common.ErrFailedToDeleteTask
	}

	return nil
}

// canViewTask checks if user can view a task
func (s *service) canViewTask(task *entity.Task, userID uuid.UUID, userRole string) bool {
	if userRole == "admin" {
		return true
	}

	if userRole == "manager" {
		return true
	}

	if task.AssignedTo != nil && *task.AssignedTo == userID {
		return true
	}

	if task.CreatedBy != nil && *task.CreatedBy == userID {
		return true
	}

	return false
}

// canUpdateTask checks if user can update a task
func (s *service) canUpdateTask(task *entity.Task, userID uuid.UUID, userRole string) bool {
	if userRole == "admin" {
		return true
	}

	if userRole == "manager" {
		return true
	}

	if task.CreatedBy != nil && *task.CreatedBy == userID {
		return true
	}

	if task.AssignedTo != nil && *task.AssignedTo == userID {
		return true
	}

	return false
}

// canDeleteTask checks if user can delete a task
func (s *service) canDeleteTask(task *entity.Task, userID uuid.UUID, userRole string) bool {
	if userRole == "admin" {
		return true
	}

	if userRole == "manager" {
		return true
	}

	if task.CreatedBy != nil && *task.CreatedBy == userID {
		return true
	}

	return false
}

// entityToResponse converts a task entity to response DTO
func (s *service) entityToResponse(task *entity.Task) *TaskResponse {
	return &TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		ProjectID:   task.ProjectID,
		CreatedBy:   task.CreatedBy,
		AssignedTo:  task.AssignedTo,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

// adjustForHolidays adjusts the date to the next business day if it's a holiday
func (s *service) adjustForHolidays(date time.Time) time.Time {
	adjustedDate := date
	for utils.CheckHoliday(adjustedDate) || adjustedDate.Weekday() == time.Saturday || adjustedDate.Weekday() == time.Sunday {
		adjustedDate = adjustedDate.AddDate(0, 0, 1)
	}
	return adjustedDate
}

package task

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mnizarzr/dot-test/common"
)

// Handler handles HTTP requests for task operations
type Handler struct {
	service Service
}

// NewHandler creates a new task handler instance
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateTask handles task creation requests
//
//	@Summary		Create a new task
//	@Description	Create a new task (must belong to a project)
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		CreateTaskRequest						true	"Task creation request"
//	@Success		201		{object}	common.BaseResponse{data=TaskResponse}	"Task created successfully"
//	@Failure		400		{object}	common.BaseResponse						"Bad request"
//	@Failure		401		{object}	common.BaseResponse						"Unauthorized"
//	@Failure		403		{object}	common.BaseResponse						"Forbidden"
//	@Failure		404		{object}	common.BaseResponse						"Project or user not found"
//	@Failure		500		{object}	common.BaseResponse						"Internal server error"
//	@Router			/api/v1/tasks [post]
func (h *Handler) CreateTask(c *gin.Context) {
	// Get user info
	userIDStr, exists := c.Get("user_id")
	if !exists {
		common.ErrorResponse(c, 401, "User ID not found")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		common.ErrorResponse(c, 401, "User role not found")
		return
	}

	// Parse request
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ValidationErrorResponse(c, err)
		return
	}

	// Create task
	task, err := h.service.CreateTask(c.Request.Context(), req, userID, userRole.(string))
	if err != nil {
		switch {
		case errors.Is(err, common.ErrProjectNotFound):
			common.ErrorResponse(c, 404, "Project not found")
			return
		case errors.Is(err, common.ErrUserNotFound):
			common.ErrorResponse(c, 404, "Assigned user not found")
			return
		case errors.Is(err, common.ErrCannotAssignTask):
			common.ErrorResponse(c, 403, "Insufficient permissions to assign tasks")
			return
		default:
			common.InternalServerErrorResponse(c, "Failed to create task")
			return
		}
	}

	common.SuccessResponse(c, task, "Task created successfully")
}

// GetTask handles get task by ID requests
//
//	@Summary		Get task by ID
//	@Description	Get a task by its ID (with permission checks)
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string									true	"Task ID"
//	@Success		200	{object}	common.BaseResponse{data=TaskResponse}	"Task retrieved successfully"
//	@Failure		400	{object}	common.BaseResponse						"Bad request"
//	@Failure		401	{object}	common.BaseResponse						"Unauthorized"
//	@Failure		403	{object}	common.BaseResponse						"Forbidden"
//	@Failure		404	{object}	common.BaseResponse						"Task not found"
//	@Failure		500	{object}	common.BaseResponse						"Internal server error"
//	@Router			/api/v1/tasks/{id} [get]
func (h *Handler) GetTask(c *gin.Context) {
	// Parse task ID
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid task ID")
		return
	}

	// Get user info
	userIDStr, exists := c.Get("user_id")
	if !exists {
		common.ErrorResponse(c, 401, "User ID not found")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		common.ErrorResponse(c, 401, "User role not found")
		return
	}

	// Get task
	task, err := h.service.GetTask(c.Request.Context(), taskID, userID, userRole.(string))
	if err != nil {
		switch {
		case errors.Is(err, common.ErrTaskNotFound):
			common.ErrorResponse(c, 404, "Task not found")
			return
		case errors.Is(err, common.ErrForbidden):
			common.ErrorResponse(c, 403, "Insufficient permissions to view task")
			return
		default:
			common.InternalServerErrorResponse(c, "Failed to retrieve task")
			return
		}
	}

	common.SuccessResponse(c, task, "Task retrieved successfully")
}

// GetTasksWithFilters handles get tasks with filters requests
//
//	@Summary		Get tasks with filters
//	@Description	Get tasks with filtering and pagination (role-based access)
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			project_id	query		string										false	"Filter by project ID"
//	@Param			assigned_to	query		string										false	"Filter by assigned user ID"
//	@Param			status		query		string										false	"Filter by status (pending, in_progress, completed)"
//	@Param			priority	query		string										false	"Filter by priority (low, medium, high)"
//	@Param			page		query		int											false	"Page number (default: 1)"
//	@Param			limit		query		int											false	"Page size (default: 10, max: 100)"
//	@Success		200			{object}	common.BaseResponse{data=TaskListResponse}	"Tasks retrieved successfully"
//	@Failure		400			{object}	common.BaseResponse							"Bad request"
//	@Failure		401			{object}	common.BaseResponse							"Unauthorized"
//	@Failure		500			{object}	common.BaseResponse							"Internal server error"
//	@Router			/api/v1/tasks [get]
func (h *Handler) GetTasksWithFilters(c *gin.Context) {
	// Get user info
	userIDStr, exists := c.Get("user_id")
	if !exists {
		common.ErrorResponse(c, 401, "User ID not found")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		common.ErrorResponse(c, 401, "User role not found")
		return
	}

	// Parse filter parameters
	var filters TaskFilterRequest
	if err := c.ShouldBindQuery(&filters); err != nil {
		common.ValidationErrorResponse(c, err)
		return
	}

	// Get tasks
	response, err := h.service.GetTasksWithFilters(c.Request.Context(), filters, userID, userRole.(string))
	if err != nil {
		common.InternalServerErrorResponse(c, "Failed to retrieve tasks")
		return
	}

	common.SuccessResponse(c, response, "Tasks retrieved successfully")
}

// UpdateTask handles task update requests
//
//	@Summary		Update task
//	@Description	Update an existing task (with permission checks)
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string									true	"Task ID"
//	@Param			request	body		UpdateTaskRequest						true	"Task update request"
//	@Success		200		{object}	common.BaseResponse{data=TaskResponse}	"Task updated successfully"
//	@Failure		400		{object}	common.BaseResponse						"Bad request"
//	@Failure		401		{object}	common.BaseResponse						"Unauthorized"
//	@Failure		403		{object}	common.BaseResponse						"Forbidden"
//	@Failure		404		{object}	common.BaseResponse						"Task not found"
//	@Failure		500		{object}	common.BaseResponse						"Internal server error"
//	@Router			/api/v1/tasks/{id} [put]
func (h *Handler) UpdateTask(c *gin.Context) {
	// Parse task ID
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid task ID")
		return
	}

	// Get user info
	userIDStr, exists := c.Get("user_id")
	if !exists {
		common.ErrorResponse(c, 401, "User ID not found")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		common.ErrorResponse(c, 401, "User role not found")
		return
	}

	// Parse request
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ValidationErrorResponse(c, err)
		return
	}

	// Update task
	task, err := h.service.UpdateTask(c.Request.Context(), taskID, req, userID, userRole.(string))
	if err != nil {
		switch {
		case errors.Is(err, common.ErrTaskNotFound):
			common.ErrorResponse(c, 404, "Task not found")
			return
		case errors.Is(err, common.ErrForbidden):
			common.ErrorResponse(c, 403, "Insufficient permissions to update task")
			return
		case errors.Is(err, common.ErrCannotAssignTask):
			common.ErrorResponse(c, 403, "Insufficient permissions to assign tasks")
			return
		case errors.Is(err, common.ErrUserNotFound):
			common.ErrorResponse(c, 404, "Assigned user not found")
			return
		default:
			common.InternalServerErrorResponse(c, "Failed to update task")
			return
		}
	}

	common.SuccessResponse(c, task, "Task updated successfully")
}

// AssignTask handles task assignment requests
//
//	@Summary		Assign task to user
//	@Description	Assign a task to a user (managers and admins only)
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string									true	"Task ID"
//	@Param			request	body		AssignTaskRequest						true	"Task assignment request"
//	@Success		200		{object}	common.BaseResponse{data=TaskResponse}	"Task assigned successfully"
//	@Failure		400		{object}	common.BaseResponse						"Bad request"
//	@Failure		401		{object}	common.BaseResponse						"Unauthorized"
//	@Failure		403		{object}	common.BaseResponse						"Forbidden"
//	@Failure		404		{object}	common.BaseResponse						"Task or user not found"
//	@Failure		500		{object}	common.BaseResponse						"Internal server error"
//	@Router			/api/v1/tasks/{id}/assign [put]
func (h *Handler) AssignTask(c *gin.Context) {
	// Parse task ID
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid task ID")
		return
	}

	// Get user info
	userIDStr, exists := c.Get("user_id")
	if !exists {
		common.ErrorResponse(c, 401, "User ID not found")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		common.ErrorResponse(c, 401, "User role not found")
		return
	}

	// Parse request
	var req AssignTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ValidationErrorResponse(c, err)
		return
	}

	// Assign task
	task, err := h.service.AssignTask(c.Request.Context(), taskID, req, userID, userRole.(string))
	if err != nil {
		switch {
		case errors.Is(err, common.ErrTaskNotFound):
			common.ErrorResponse(c, 404, "Task not found")
			return
		case errors.Is(err, common.ErrCannotAssignTask):
			common.ErrorResponse(c, 403, "Insufficient permissions to assign tasks")
			return
		case errors.Is(err, common.ErrUserNotFound):
			common.ErrorResponse(c, 404, "Assigned user not found")
			return
		default:
			common.InternalServerErrorResponse(c, "Failed to assign task")
			return
		}
	}

	common.SuccessResponse(c, task, "Task assigned successfully")
}

// DeleteTask handles task deletion requests
//
//	@Summary		Delete task
//	@Description	Delete an existing task (creator, manager, or admin only)
//	@Tags			Task
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string				true	"Task ID"
//	@Success		200	{object}	common.BaseResponse	"Task deleted successfully"
//	@Failure		400	{object}	common.BaseResponse	"Bad request"
//	@Failure		401	{object}	common.BaseResponse	"Unauthorized"
//	@Failure		403	{object}	common.BaseResponse	"Forbidden"
//	@Failure		404	{object}	common.BaseResponse	"Task not found"
//	@Failure		500	{object}	common.BaseResponse	"Internal server error"
//	@Router			/api/v1/tasks/{id} [delete]
func (h *Handler) DeleteTask(c *gin.Context) {
	// Parse task ID
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid task ID")
		return
	}

	// Get user info
	userIDStr, exists := c.Get("user_id")
	if !exists {
		common.ErrorResponse(c, 401, "User ID not found")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		common.ErrorResponse(c, 401, "User role not found")
		return
	}

	// Delete task
	err = h.service.DeleteTask(c.Request.Context(), taskID, userID, userRole.(string))
	if err != nil {
		switch {
		case errors.Is(err, common.ErrTaskNotFound):
			common.ErrorResponse(c, 404, "Task not found")
			return
		case errors.Is(err, common.ErrForbidden):
			common.ErrorResponse(c, 403, "Insufficient permissions to delete task")
			return
		default:
			common.InternalServerErrorResponse(c, "Failed to delete task")
			return
		}
	}

	common.SuccessResponse(c, nil, "Task deleted successfully")
}

package project

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mnizarzr/dot-test/common"
)

// Handler handles HTTP requests for project operations
type Handler struct {
	service Service
}

// NewHandler creates a new project handler instance
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateProject handles project creation requests
//
//	@Summary		Create a new project
//	@Description	Create a new project (only managers and admins)
//	@Tags			Project
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			payload	body		CreateProjectRequest						true	"Project creation request"
//	@Success		201		{object}	common.BaseResponse{data=ProjectResponse}	"Project created successfully"
//	@Failure		400		{object}	common.BaseResponse							"Bad request"
//	@Failure		401		{object}	common.BaseResponse							"Unauthorized"
//	@Failure		403		{object}	common.BaseResponse							"Forbidden"
//	@Failure		409		{object}	common.BaseResponse							"Project name already exists"
//	@Failure		500		{object}	common.BaseResponse							"Internal server error"
//	@Router			/api/v1/projects [post]
func (h *Handler) CreateProject(c *gin.Context) {
	// Check role permissions
	userRole, exists := c.Get("user_role")
	if !exists {
		common.ErrorResponse(c, 401, "User role not found")
		return
	}

	role := userRole.(string)
	if role != "admin" && role != "manager" {
		common.ErrorResponse(c, 403, "Only managers and admins can create projects")
		return
	}

	// Get user ID
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

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ValidationErrorResponse(c, err)
		return
	}

	project, err := h.service.CreateProject(c.Request.Context(), req, userID)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrProjectNameAlreadyExists):
			common.ErrorResponse(c, 409, "Project name already exists")
			return
		default:
			common.InternalServerErrorResponse(c, "Failed to create project")
			return
		}
	}

	common.SuccessResponse(c, project, "Project created successfully")
}

// GetProject handles get project by ID requests
//
//	@Summary		Get project by ID
//	@Description	Get a project by its ID
//	@Tags			Project
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string										true	"Project ID"
//	@Success		200	{object}	common.BaseResponse{data=ProjectResponse}	"Project retrieved successfully"
//	@Failure		400	{object}	common.BaseResponse							"Bad request"
//	@Failure		401	{object}	common.BaseResponse							"Unauthorized"
//	@Failure		404	{object}	common.BaseResponse							"Project not found"
//	@Failure		500	{object}	common.BaseResponse							"Internal server error"
//	@Router			/api/v1/projects/{id} [get]
func (h *Handler) GetProject(c *gin.Context) {
	// Parse project ID
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid project ID")
		return
	}

	project, err := h.service.GetProject(c.Request.Context(), projectID)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrProjectNotFound):
			common.ErrorResponse(c, 404, "Project not found")
			return
		default:
			common.InternalServerErrorResponse(c, "Failed to retrieve project")
			return
		}
	}

	common.SuccessResponse(c, project, "Project retrieved successfully")
}

// GetAllProjects handles get all projects requests
//
//	@Summary		Get all projects
//	@Description	Get all projects with pagination
//	@Tags			Project
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			page	query		int												false	"Page number (default: 1)"
//	@Param			limit	query		int												false	"Page size (default: 10, max: 100)"
//	@Success		200		{object}	common.BaseResponse{data=ProjectListResponse}	"Projects retrieved successfully"
//	@Failure		400		{object}	common.BaseResponse								"Bad request"
//	@Failure		401		{object}	common.BaseResponse								"Unauthorized"
//	@Failure		500		{object}	common.BaseResponse								"Internal server error"
//	@Router			/api/v1/projects [get]
func (h *Handler) GetAllProjects(c *gin.Context) {
	var pagination PaginationRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		common.ValidationErrorResponse(c, err)
		return
	}

	response, err := h.service.GetAllProjects(c.Request.Context(), pagination)
	if err != nil {
		common.InternalServerErrorResponse(c, "Failed to retrieve projects")
		return
	}

	common.SuccessResponse(c, response, "Projects retrieved successfully")
}

// UpdateProject handles project update requests
//
//	@Summary		Update project
//	@Description	Update an existing project
//	@Tags			Project
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path		string										true	"Project ID"
//	@Param			request	body		UpdateProjectRequest						true	"Project update request"
//	@Success		200		{object}	common.BaseResponse{data=ProjectResponse}	"Project updated successfully"
//	@Failure		400		{object}	common.BaseResponse							"Bad request"
//	@Failure		401		{object}	common.BaseResponse							"Unauthorized"
//	@Failure		403		{object}	common.BaseResponse							"Forbidden"
//	@Failure		404		{object}	common.BaseResponse							"Project not found"
//	@Failure		409		{object}	common.BaseResponse							"Project name already exists"
//	@Failure		500		{object}	common.BaseResponse							"Internal server error"
//	@Router			/api/v1/projects/{id} [put]
func (h *Handler) UpdateProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid project ID")
		return
	}

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
	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ValidationErrorResponse(c, err)
		return
	}

	project, err := h.service.UpdateProject(c.Request.Context(), projectID, req, userID, userRole.(string))
	if err != nil {
		switch {
		case errors.Is(err, common.ErrProjectNotFound):
			common.ErrorResponse(c, 404, "Project not found")
			return
		case errors.Is(err, common.ErrForbidden):
			common.ErrorResponse(c, 403, "Insufficient permissions to update project")
			return
		case errors.Is(err, common.ErrProjectNameAlreadyExists):
			common.ErrorResponse(c, 409, "Project name already exists")
			return
		default:
			common.InternalServerErrorResponse(c, "Failed to update project")
			return
		}
	}

	common.SuccessResponse(c, project, "Project updated successfully")
}

// DeleteProject handles project deletion requests
//
//	@Summary		Delete project
//	@Description	Delete an existing project and all its tasks
//	@Tags			Project
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		string				true	"Project ID"
//	@Success		200	{object}	common.BaseResponse	"Project deleted successfully"
//	@Failure		400	{object}	common.BaseResponse	"Bad request"
//	@Failure		401	{object}	common.BaseResponse	"Unauthorized"
//	@Failure		403	{object}	common.BaseResponse	"Forbidden"
//	@Failure		404	{object}	common.BaseResponse	"Project not found"
//	@Failure		500	{object}	common.BaseResponse	"Internal server error"
//	@Router			/api/v1/projects/{id} [delete]
func (h *Handler) DeleteProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid project ID")
		return
	}

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

	err = h.service.DeleteProject(c.Request.Context(), projectID, userID, userRole.(string))
	if err != nil {
		switch {
		case errors.Is(err, common.ErrProjectNotFound):
			common.ErrorResponse(c, 404, "Project not found")
			return
		case errors.Is(err, common.ErrForbidden):
			common.ErrorResponse(c, 403, "Insufficient permissions to delete project")
			return
		default:
			common.InternalServerErrorResponse(c, "Failed to delete project")
			return
		}
	}

	common.SuccessResponse(c, nil, "Project deleted successfully")
}

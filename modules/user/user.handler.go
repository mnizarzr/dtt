package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mnizarzr/dot-test/common"
)

// Handler handles HTTP requests for user operations
type Handler struct {
	service Service
}

// NewHandler creates a new user handler instance
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetProfile handles get user profile requests
//
//	@Summary		Get user profile
//	@Description	Get the current user's profile information
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	common.BaseResponse{data=ProfileResponse}	"User profile retrieved successfully"
//	@Failure		401	{object}	common.BaseResponse							"Unauthorized"
//	@Failure		404	{object}	common.BaseResponse							"User not found"
//	@Failure		500	{object}	common.BaseResponse							"Internal server error"
//	@Router			/api/v1/user/me [get]
func (h *Handler) GetProfile(c *gin.Context) {
	// Get user ID from JWT middleware context
	userIDStr, exists := c.Get("user_id")
	if !exists {
		common.ErrorResponse(c, 401, "User ID not found in token")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		common.ErrorResponse(c, 400, "Invalid user ID format")
		return
	}

	// Get user profile
	profile, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrUserNotFound):
			common.ErrorResponse(c, 404, "User not found")
			return
		default:
			common.InternalServerErrorResponse(c, "Failed to retrieve profile")
			return
		}
	}

	common.SuccessResponse(c, profile, "Profile retrieved successfully")
}

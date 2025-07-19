package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mnizarzr/dot-test/common"
)

// Handler handles HTTP requests for auth operations
type Handler struct {
	service Service
}

// NewHandler creates a new auth handler instance
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Register handles user registration requests
//
//	@Summary		Register a new user
//	@Description	Register a new user account with email and password. Admins can specify role.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterRequest								true	"User registration data"
//	@Success		201		{object}	common.BaseResponse{data=RegisterResponse}	"User registered successfully"
//	@Failure		400		{object}	common.BaseResponse							"Bad request - validation errors"
//	@Failure		401		{object}	common.BaseResponse							"Unauthorized - invalid token"
//	@Failure		403		{object}	common.BaseResponse							"Forbidden - insufficient permissions"
//	@Failure		409		{object}	common.BaseResponse							"Conflict - email already exists"
//	@Failure		422		{object}	common.BaseResponse{data=ValidationErrors}	"Validation failed"
//	@Failure		500		{object}	common.BaseResponse							"Internal server error"
//	@Router			/api/v1/auth/register [post]
//	@Security		BearerAuth
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	// Bind and validate JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := h.extractValidationErrors(err)
		if len(validationErrors.Errors) > 0 {
			common.ValidationErrorResponse(c, validationErrors)
			return
		}
		common.BadRequestResponse(c, "Invalid request format")
		return
	}

	requestingUserRole := ""
	if userRole, exists := c.Get("user_role"); exists {
		requestingUserRole = userRole.(string)
	}

	response, err := h.service.Register(c.Request.Context(), req, requestingUserRole)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrEmailAlreadyRegistered):
			common.ConflictResponse(c, "Email address is already registered")
			return
		case errors.Is(err, common.ErrInvalidNameFormat):
			common.BadRequestResponse(c, err.Error())
			return
		case errors.Is(err, common.ErrInvalidEmailFormat):
			common.BadRequestResponse(c, err.Error())
			return
		case errors.Is(err, common.ErrInvalidPasswordFormat):
			common.BadRequestResponse(c, err.Error())
			return
		case err.Error() == "insufficient permissions to assign role":
			common.ErrorResponse(c, 403, "Insufficient permissions to assign role")
			return
		default:
			// Check for invalid role errors
			if len(err.Error()) > 13 && err.Error()[:13] == "invalid role:" {
				common.BadRequestResponse(c, err.Error())
				return
			}
			common.InternalServerErrorResponse(c, "Failed to register user")
			return
		}
	}

	common.CreatedResponse(c, response, "User registered successfully")
}

// Login handles user login requests
//
//	@Summary		User login
//	@Description	Authenticate user with email and password
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest								true	"User login credentials"
//	@Success		200		{object}	common.BaseResponse{data=LoginResponse}		"Login successful"
//	@Failure		400		{object}	common.BaseResponse							"Bad request - validation errors"
//	@Failure		401		{object}	common.BaseResponse							"Unauthorized - invalid credentials"
//	@Failure		422		{object}	common.BaseResponse{data=ValidationErrors}	"Validation failed"
//	@Failure		500		{object}	common.BaseResponse							"Internal server error"
//	@Router			/api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		validationErrors := h.extractValidationErrors(err)
		if len(validationErrors.Errors) > 0 {
			common.ValidationErrorResponse(c, validationErrors)
			return
		}
		common.BadRequestResponse(c, "Invalid request format")
		return
	}

	response, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidCredentials):
			common.ErrorResponse(c, 401, "Invalid email or password")
			return
		case errors.Is(err, common.ErrInvalidEmailFormat):
			common.BadRequestResponse(c, err.Error())
			return
		case errors.Is(err, common.ErrInvalidPasswordFormat):
			common.BadRequestResponse(c, err.Error())
			return
		default:
			common.InternalServerErrorResponse(c, "Failed to login user")
			return
		}
	}

	common.SuccessResponse(c, response, "Login successful")
}

// extractValidationErrors extracts validation errors from binding errors
func (h *Handler) extractValidationErrors(err error) ValidationErrors {
	var validationErrors ValidationErrors

	validationErrors.Errors = append(validationErrors.Errors, ValidationError{
		Field:   "general",
		Message: err.Error(),
	})

	return validationErrors
}

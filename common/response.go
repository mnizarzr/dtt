package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BaseResponse represents the standard API response structure
type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SuccessResponse creates a success response
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	if message == "" {
		message = "Success"
	}

	response := BaseResponse{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	}

	c.JSON(http.StatusOK, response)
}

// CreatedResponse creates a resource created response
func CreatedResponse(c *gin.Context, data interface{}, message string) {
	if message == "" {
		message = "Resource created successfully"
	}

	response := BaseResponse{
		Code:    http.StatusCreated,
		Message: message,
		Data:    data,
	}

	c.JSON(http.StatusCreated, response)
}

// ErrorResponse creates an error response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	response := BaseResponse{
		Code:    statusCode,
		Message: message,
		Data:    nil,
	}

	c.JSON(statusCode, response)
}

// BadRequestResponse creates a bad request error response
func BadRequestResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Bad request"
	}
	ErrorResponse(c, http.StatusBadRequest, message)
}

// ValidationErrorResponse creates a validation error response
func ValidationErrorResponse(c *gin.Context, errors interface{}) {
	response := BaseResponse{
		Code:    http.StatusUnprocessableEntity,
		Message: "Validation failed",
		Data:    errors,
	}

	c.JSON(http.StatusUnprocessableEntity, response)
}

// InternalServerErrorResponse creates an internal server error response
func InternalServerErrorResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	ErrorResponse(c, http.StatusInternalServerError, message)
}

// ConflictResponse creates a conflict error response
func ConflictResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Resource already exists"
	}
	ErrorResponse(c, http.StatusConflict, message)
}

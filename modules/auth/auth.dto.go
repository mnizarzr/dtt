package auth

import (
	"time"

	"github.com/google/uuid"
)

// RegisterRequest represents the request payload for user registration
type RegisterRequest struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required" example:"SecurePass123"`
	Role     string `json:"role,omitempty" example:"user"` // Optional, only allowed for admin users
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required" example:"SecurePass123"`
}

// UserResponse represents the user data in API responses
type UserResponse struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Role      string    `json:"role" example:"user"`
	CreatedAt time.Time `json:"created_at" example:"2025-07-19T10:30:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-07-19T10:30:00Z"`
}

// RegisterResponse represents the response for user registration
type RegisterResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message" example:"Registration successful. Welcome email has been sent."`
}

// LoginResponse represents the response for user login
type LoginResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string       `json:"token_type" example:"Bearer"`
	ExpiresIn   int64        `json:"expires_in" example:"3600"`
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"email format is invalid"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

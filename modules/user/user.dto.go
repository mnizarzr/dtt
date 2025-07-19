package user

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse represents the user data in API responses
type UserResponse struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Role      string    `json:"role" example:"user"`
	CreatedAt time.Time `json:"created_at" example:"2025-07-19T10:30:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-07-19T10:30:00Z"`
}

// ProfileResponse represents the response for user profile
type ProfileResponse struct {
	User UserResponse `json:"user"`
}

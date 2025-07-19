package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/mnizarzr/dot-test/common"
	"github.com/mnizarzr/dot-test/entity"
)

// Service defines the interface for user business logic
type Service interface {
	GetProfile(ctx context.Context, id uuid.UUID) (*ProfileResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponse, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new user service instance
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// GetProfile retrieves a user profile by ID
func (s *service) GetProfile(ctx context.Context, id uuid.UUID) (*ProfileResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, common.ErrFailedToRetrieveUser
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}

	userResponse := s.entityToResponse(user)
	response := &ProfileResponse{
		User: *userResponse,
	}

	return response, nil
}

// GetUserByEmail retrieves a user by email
func (s *service) GetUserByEmail(ctx context.Context, email string) (*UserResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, common.ErrFailedToRetrieveUser
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}

	return s.entityToResponse(user), nil
}

// GetUserByID retrieves a user by ID
func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, common.ErrFailedToRetrieveUser
	}
	if user == nil {
		return nil, common.ErrUserNotFound
	}

	return s.entityToResponse(user), nil
}

// entityToResponse converts a user entity to response DTO
func (s *service) entityToResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

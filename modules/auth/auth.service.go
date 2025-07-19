package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/mnizarzr/dot-test/common"
	"github.com/mnizarzr/dot-test/entity"
	"github.com/mnizarzr/dot-test/jobs"
	"github.com/mnizarzr/dot-test/modules/user"
	"github.com/mnizarzr/dot-test/utils"
)

// Service defines the interface for auth business logic
type Service interface {
	Register(ctx context.Context, req RegisterRequest, requestingUserRole string) (*RegisterResponse, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
}

// service implements the Service interface
type service struct {
	userRepo  user.Repository
	jobClient *asynq.Client
	jwtSecret string
}

// NewService creates a new auth service instance
func NewService(userRepo user.Repository, jobClient *asynq.Client, jwtSecret string) Service {
	return &service{
		userRepo:  userRepo,
		jobClient: jobClient,
		jwtSecret: jwtSecret,
	}
}

// Register handles user registration business logic
func (s *service) Register(ctx context.Context, req RegisterRequest, requestingUserRole string) (*RegisterResponse, error) {
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	userRole, err := s.determineUserRole(req.Role, requestingUserRole)
	if err != nil {
		return nil, err
	}

	exists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, common.ErrFailedToCheckEmail
	}
	if exists {
		return nil, common.ErrEmailAlreadyRegistered
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, common.ErrFailedToHashPassword
	}

	now := time.Now()
	userEntity := &entity.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         userRole,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, userEntity); err != nil {
		return nil, common.ErrFailedToCreateUser
	}

	if err := s.enqueueWelcomeEmail(ctx, userEntity.Email, userEntity.Name, userEntity.Role, req.Password); err != nil {
		log.Printf("Failed to enqueue welcome email for user %s: %v", userEntity.Email, err)
	}

	userResponse := UserResponse{
		ID:        userEntity.ID,
		Name:      userEntity.Name,
		Email:     userEntity.Email,
		Role:      userEntity.Role,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}

	response := &RegisterResponse{
		User:    userResponse,
		Message: "Registration successful. Welcome email has been sent.",
	}

	return response, nil
}

// Login handles user login business logic
func (s *service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if err := s.validateLoginRequest(req); err != nil {
		return nil, err
	}

	userEntity, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, common.ErrFailedToRetrieveUser
	}
	if userEntity == nil {
		return nil, common.ErrInvalidCredentials
	}

	if !utils.CheckPasswordHash(req.Password, userEntity.PasswordHash) {
		return nil, common.ErrInvalidCredentials
	}

	token, expiresIn, err := utils.GenerateJWT(userEntity.ID.String(), userEntity.Email, userEntity.Role, s.jwtSecret)
	if err != nil {
		return nil, common.ErrFailedToGenerateToken
	}

	userResponse := UserResponse{
		ID:        userEntity.ID,
		Name:      userEntity.Name,
		Email:     userEntity.Email,
		Role:      userEntity.Role,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}

	response := &LoginResponse{
		User:        userResponse,
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}

	return response, nil
}

// determineUserRole determines the user role based on request and permissions
func (s *service) determineUserRole(requestedRole, requestingUserRole string) (string, error) {
	if requestedRole == "" {
		return "user", nil
	}

	validRoles := map[string]bool{
		"user":    true,
		"manager": true,
		"admin":   true,
	}

	if !validRoles[requestedRole] {
		return "", fmt.Errorf("invalid role: %s", requestedRole)
	}

	if requestedRole == "user" {
		return "user", nil
	}

	if requestingUserRole != "admin" {
		return "", fmt.Errorf("insufficient permissions to assign role")
	}

	return requestedRole, nil
}

// validateRegisterRequest validates the registration request
func (s *service) validateRegisterRequest(req RegisterRequest) error {
	if !utils.IsValidName(req.Name) {
		return common.ErrInvalidNameFormat
	}
	if !utils.IsValidEmail(req.Email) {
		return common.ErrInvalidEmailFormat
	}
	if !utils.IsValidPassword(req.Password) {
		return common.ErrInvalidPasswordFormat
	}

	return nil
}

// validateLoginRequest validates the login request
func (s *service) validateLoginRequest(req LoginRequest) error {
	if !utils.IsValidEmail(req.Email) {
		return common.ErrInvalidEmailFormat
	}

	if len(req.Password) == 0 {
		return common.ErrInvalidPasswordFormat
	}

	return nil
}

// enqueueWelcomeEmail adds a welcome email job to the queue
func (s *service) enqueueWelcomeEmail(ctx context.Context, email, name, role, password string) error {
	task, err := jobs.NewWelcomeEmailTask(email, name, role, password)
	if err != nil {
		return err
	}

	// Enqueue with high priority for welcome emails
	_, err = s.jobClient.Enqueue(task, asynq.Queue("default"), asynq.MaxRetry(3))
	return err
}

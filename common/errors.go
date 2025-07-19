package common

import "errors"

var ErrForbidden = errors.New("forbidden")

// User-related errors
var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrInvalidEmailFormat     = errors.New("invalid email format")
	ErrInvalidNameFormat      = errors.New("name must be between 2-100 characters and contain only letters, spaces, hyphens, and apostrophes")
	ErrInvalidPasswordFormat  = errors.New("password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one digit")
	ErrUserNotFound           = errors.New("user not found")
	ErrFailedToCreateUser     = errors.New("failed to create user")
	ErrFailedToRetrieveUser   = errors.New("failed to retrieve user")
	ErrFailedToCheckEmail     = errors.New("failed to check email existence")
	ErrFailedToHashPassword   = errors.New("failed to hash password")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrFailedToGenerateToken  = errors.New("failed to generate authentication token")
)

// Project-related errors
var (
	ErrProjectNotFound          = errors.New("project not found")
	ErrProjectNameAlreadyExists = errors.New("project name already exists")
	ErrFailedToCreateProject    = errors.New("failed to create project")
	ErrFailedToRetrieveProject  = errors.New("failed to retrieve project")
	ErrFailedToRetrieveProjects = errors.New("failed to retrieve projects")
	ErrFailedToUpdateProject    = errors.New("failed to update project")
	ErrFailedToDeleteProject    = errors.New("failed to delete project")
	ErrFailedToCheckProject     = errors.New("failed to check project existence")
)

// Task-related errors
var (
	ErrTaskNotFound          = errors.New("task not found")
	ErrFailedToCreateTask    = errors.New("failed to create task")
	ErrFailedToRetrieveTask  = errors.New("failed to retrieve task")
	ErrFailedToRetrieveTasks = errors.New("failed to retrieve tasks")
	ErrFailedToUpdateTask    = errors.New("failed to update task")
	ErrFailedToDeleteTask    = errors.New("failed to delete task")
	ErrCannotAssignTask      = errors.New("insufficient permissions to assign tasks")
)

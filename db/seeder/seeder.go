package seeder

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mnizarzr/dot-test/entity"
	"github.com/mnizarzr/dot-test/utils"
	"gorm.io/gorm"
)

// Seeder handles database seeding operations
type Seeder struct {
	db *gorm.DB
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{
		db: db,
	}
}

// CreateAdminUser creates an admin user if it doesn't exist
func (s *Seeder) CreateAdminUser(ctx context.Context, name, email, password string) (*entity.User, error) {
	// Check if admin user already exists
	var existingUser entity.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&existingUser).Error; err == nil {
		return &existingUser, fmt.Errorf("user with email %s already exists", email)
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}

	// Validate input
	if !utils.IsValidName(name) {
		return nil, fmt.Errorf("invalid name format")
	}
	if !utils.IsValidEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}
	if !utils.IsValidPassword(password) {
		return nil, fmt.Errorf("invalid password format")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create admin user
	now := time.Now()
	adminUser := &entity.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         "admin",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Save to database
	if err := s.db.WithContext(ctx).Create(adminUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}

	return adminUser, nil
}

// CreateDefaultUsers creates default users for testing
func (s *Seeder) CreateDefaultUsers(ctx context.Context) error {
	defaultUsers := []struct {
		Name     string
		Email    string
		Password string
		Role     string
	}{
		{
			Name:     "System Admin",
			Email:    "admin@dottest.com",
			Password: "AdminPass123!",
			Role:     "admin",
		},
		{
			Name:     "Project Manager",
			Email:    "manager@dottest.com",
			Password: "ManagerPass123!",
			Role:     "manager",
		},
		{
			Name:     "Test User",
			Email:    "user@dottest.com",
			Password: "UserPass123!",
			Role:     "user",
		},
	}

	for _, userData := range defaultUsers {
		// Check if user already exists
		var existingUser entity.User
		if err := s.db.WithContext(ctx).Where("email = ?", userData.Email).First(&existingUser).Error; err == nil {
			fmt.Printf("User %s already exists, skipping...\n", userData.Email)
			continue
		} else if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to check if user %s exists: %w", userData.Email, err)
		}

		// Hash password
		hashedPassword, err := utils.HashPassword(userData.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", userData.Email, err)
		}

		// Create user
		now := time.Now()
		user := &entity.User{
			ID:           uuid.New(),
			Name:         userData.Name,
			Email:        userData.Email,
			PasswordHash: hashedPassword,
			Role:         userData.Role,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		// Save to database
		if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
			return fmt.Errorf("failed to create user %s: %w", userData.Email, err)
		}

		fmt.Printf("✅ Created %s user: %s\n", userData.Role, userData.Email)
	}

	return nil
}

// SeedProjects creates sample projects
func (s *Seeder) SeedProjects(ctx context.Context) error {
	// Get a manager or admin user to be the creator
	var managerUser entity.User
	if err := s.db.WithContext(ctx).Where("role IN ?", []string{"manager", "admin"}).First(&managerUser).Error; err != nil {
		return fmt.Errorf("no manager or admin user found to create projects: %w", err)
	}

	sampleProjects := []struct {
		Name        string
		Description string
	}{
		{
			Name:        "Website Redesign",
			Description: "Complete redesign of the company website with modern UI/UX",
		},
		{
			Name:        "Mobile App Development",
			Description: "Development of iOS and Android mobile application",
		},
		{
			Name:        "Database Migration",
			Description: "Migration of legacy database to PostgreSQL",
		},
	}

	for _, projectData := range sampleProjects {
		// Check if project already exists
		var existingProject entity.Project
		if err := s.db.WithContext(ctx).Where("name = ?", projectData.Name).First(&existingProject).Error; err == nil {
			fmt.Printf("Project %s already exists, skipping...\n", projectData.Name)
			continue
		} else if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("failed to check if project %s exists: %w", projectData.Name, err)
		}

		// Create project
		now := time.Now()
		project := &entity.Project{
			ID:          uuid.New(),
			Name:        projectData.Name,
			Description: projectData.Description,
			CreatedBy:   &managerUser.ID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// Save to database
		if err := s.db.WithContext(ctx).Create(project).Error; err != nil {
			return fmt.Errorf("failed to create project %s: %w", projectData.Name, err)
		}

		fmt.Printf("✅ Created project: %s\n", projectData.Name)
	}

	return nil
}

// SeedAll runs all seeders
func (s *Seeder) SeedAll(ctx context.Context) error {
	fmt.Println("🌱 Starting database seeding...")

	// Seed default users
	if err := s.CreateDefaultUsers(ctx); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	// Seed projects
	if err := s.SeedProjects(ctx); err != nil {
		return fmt.Errorf("failed to seed projects: %w", err)
	}

	fmt.Println("🎉 Database seeding completed successfully!")
	return nil
}

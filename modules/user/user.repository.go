package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mnizarzr/dot-test/db"
	"github.com/mnizarzr/dot-test/entity"
	"gorm.io/gorm"
)

// Repository defines the interface for user data operations
type Repository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	EmailExists(ctx context.Context, email string) (bool, error)
}

// repository implements the Repository interface
type repository struct {
	db    *gorm.DB
	cache *db.RedisClient
}

// NewRepository creates a new user repository instance
func NewRepository(database *gorm.DB, cache *db.RedisClient) Repository {
	return &repository{
		db:    database,
		cache: cache,
	}
}

// Create creates a new user in the database
func (r *repository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByEmail retrieves a user by email address with read-through cache
func (r *repository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	cacheKey := fmt.Sprintf("user:email:%s", email)

	// Try to get from cache first
	var user entity.User
	err := r.cache.Get(ctx, cacheKey, &user)
	if err == nil && user.ID != uuid.Nil {
		return &user, nil
	}

	// Cache miss or error, query database
	err = r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Store in cache for 1 hour
	_ = r.cache.Set(ctx, cacheKey, &user, time.Hour)

	return &user, nil
}

// GetByID retrieves a user by ID with read-through cache
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	cacheKey := fmt.Sprintf("user:id:%s", id.String())

	// Try to get from cache first
	var user entity.User
	err := r.cache.Get(ctx, cacheKey, &user)
	if err == nil && user.ID != uuid.Nil {
		return &user, nil
	}

	// Cache miss or error, query database
	err = r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Store in cache for 1 hour
	_ = r.cache.Set(ctx, cacheKey, &user, time.Hour)

	return &user, nil
}

// Update updates an existing user
func (r *repository) Update(ctx context.Context, user *entity.User) error {
	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return err
	}

	// Invalidate cache after update
	cacheKeys := []string{
		fmt.Sprintf("user:id:%s", user.ID.String()),
		fmt.Sprintf("user:email:%s", user.Email),
	}
	_ = r.cache.Delete(ctx, cacheKeys...)

	return nil
}

// Delete soft deletes a user by ID
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	// Get user first to invalidate email cache
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return nil // User doesn't exist
	}

	err = r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
	if err != nil {
		return err
	}

	// Invalidate cache after delete
	cacheKeys := []string{
		fmt.Sprintf("user:id:%s", id.String()),
		fmt.Sprintf("user:email:%s", user.Email),
	}
	_ = r.cache.Delete(ctx, cacheKeys...)

	return nil
}

// EmailExists checks if an email address is already registered
func (r *repository) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

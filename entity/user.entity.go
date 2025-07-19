package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	userTableName = "users"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (*User) TableName() string {
	return userTableName
}

func (u *User) AfterCreate(tx *gorm.DB) error {
	return CreateAuditLog(tx, AuditActionCreate, userTableName, u.ID, u, nil)
}

func (u *User) AfterUpdate(tx *gorm.DB) error {
	return CreateAuditLog(tx, AuditActionUpdate, userTableName, u.ID, u, nil)
}

func (u *User) AfterDelete(tx *gorm.DB) error {
	return CreateAuditLog(tx, AuditActionDelete, userTableName, u.ID, nil, u)
}

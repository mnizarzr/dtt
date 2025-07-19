package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	projectTableName = "projects"
)

type Project struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedBy   *uuid.UUID `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (*Project) TableName() string {
	return projectTableName
}

func (p *Project) AfterCreate(tx *gorm.DB) error {
	return CreateAuditLog(tx, AuditActionCreate, projectTableName, p.ID, p, nil)
}

func (p *Project) AfterUpdate(tx *gorm.DB) error {
	return CreateAuditLog(tx, AuditActionUpdate, projectTableName, p.ID, p, nil)
}

func (p *Project) AfterDelete(tx *gorm.DB) error {
	return CreateAuditLog(tx, AuditActionDelete, projectTableName, p.ID, nil, p)
}

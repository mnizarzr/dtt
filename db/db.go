package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresGormDb(uri string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}

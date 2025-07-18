package entity

import "github.com/google/uuid"

const (
	userTableName = "users"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	UserName string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

func (*User) TableName() string {
	return userTableName
}

package models

import "time"

type UserModel struct {
	ID           int64
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdateAt     *time.Time
}

type CreateUserModel struct {
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdateAt     *time.Time
}

type UpdateUserModel struct {
	FirstName string
	LastName  string
	UpdatedAt time.Time
}

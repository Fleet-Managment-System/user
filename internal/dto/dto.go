package dto

import "time"

type CreateUserDto struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type UpdateUserDto struct {
	FirstName string
	LastName  string
}

type GetUserDto struct {
	Id           int64
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdateAt     *time.Time
}

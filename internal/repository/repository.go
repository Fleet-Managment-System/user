package repository

import (
	"context"
	"user/internal/domain/models"
)

type Pagination struct {
	Limit  uint64
	Offset uint64
}

type IUserRepository interface {
	Create(ctx context.Context, user *models.CreateUserModel) (int64, error)
	Update(ctx context.Context, userId int64, user *models.UpdateUserModel) error
	GetOne(ctx context.Context, userId int64) (*models.UserModel, error)
	GetAll(ctx context.Context, pagination *Pagination) ([]*models.UserModel, error)
	Delete(ctx context.Context, userId int64) error
}

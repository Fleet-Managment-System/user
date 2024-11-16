package service

import (
	"context"
	"time"
	"user/internal/converters"
	"user/internal/domain/models"
	"user/internal/dto"

	rootRepository "user/internal/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HashPassword(password string) (string, error) {
	return password, nil
}

type IUserService interface {
	Create(ctx context.Context, createUserDto *dto.CreateUserDto) (int64, error)
	Update(ctx context.Context, userId int64, user *dto.UpdateUserDto) error
	GetOne(ctx context.Context, userId int64) (*dto.GetUserDto, error)
	GetAll(ctx context.Context) ([]*dto.GetUserDto, error)
	Delete(ctx context.Context, userId int64) error
}

type UserService struct {
	repository rootRepository.IUserRepository
}

func (userService *UserService) Create(ctx context.Context, createUserDto *dto.CreateUserDto) (int64, error) {
	passwordHash, err := HashPassword(createUserDto.Password)
	if err != nil {
		return -1, status.Errorf(codes.Internal, "Validation error: %v \n", err)
	}

	id, err := userService.repository.Create(ctx, &models.CreateUserModel{
		FirstName:    createUserDto.FirstName,
		LastName:     createUserDto.LastName,
		Email:        createUserDto.Email,
		PasswordHash: passwordHash,
	})

	return id, err
}

func (userService *UserService) Update(ctx context.Context, userId int64, updateUserDto *dto.UpdateUserDto) error {
	// Add transactions

	user, err := userService.repository.GetOne(ctx, userId)
	if err != nil {
		return err
	}
	if user == nil {
		return status.Errorf(codes.InvalidArgument, "User not found\n")
	}

	return userService.repository.Update(ctx, userId, &models.UpdateUserModel{
		FirstName: updateUserDto.FirstName,
		LastName:  updateUserDto.LastName,
		UpdatedAt: time.Now(),
	})
}

func (userService *UserService) Delete(ctx context.Context, userId int64) error {
	return userService.repository.Delete(ctx, userId)
}

func (userService *UserService) GetOne(ctx context.Context, userId int64) (*dto.GetUserDto, error) {
	user, err := userService.repository.GetOne(ctx, userId)
	if err != nil {
		return nil, err
	}

	return converters.UserModelToGetUserDto(user), nil
}

func (userService *UserService) GetAll(ctx context.Context) ([]*dto.GetUserDto, error) {
	userModels, err := userService.repository.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	getUserDtos := []*dto.GetUserDto{}
	for _, userModel := range userModels {
		getUserDtos = append(getUserDtos, converters.UserModelToGetUserDto(userModel))
	}

	return getUserDtos, nil
}

func NewUserService(userRepository rootRepository.IUserRepository) *UserService {
	return &UserService{
		repository: userRepository,
	}
}

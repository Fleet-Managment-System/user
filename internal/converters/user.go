package converters

import (
	"user/internal/domain/models"
	"user/internal/dto"
	desc "user/pkg/user_v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateUserRequestToCreateUserDto(createUserRequest *desc.CreateUserRequest) *dto.CreateUserDto {
	return &dto.CreateUserDto{
		FirstName: createUserRequest.User.GetFirstname(),
		LastName:  createUserRequest.User.GetLastname(),
		Email:     createUserRequest.User.GetEmail(),
		Password:  createUserRequest.User.GetPassword(),
	}
}

func UserModelToGetUserDto(userModel *models.UserModel) *dto.GetUserDto {
	return &dto.GetUserDto{
		Id:           userModel.ID,
		FirstName:    userModel.FirstName,
		LastName:     userModel.LastName,
		Email:        userModel.Email,
		PasswordHash: userModel.PasswordHash,
		CreatedAt:    userModel.CreatedAt,
		UpdateAt:     userModel.UpdateAt,
	}
}

func GetUserDtoToUserRequest(getUserDto *dto.GetUserDto) *desc.User {
	dist := &desc.User{
		Id:           getUserDto.Id,
		Firstname:    getUserDto.FirstName,
		Lastname:     getUserDto.LastName,
		Email:        getUserDto.Email,
		PasswordHash: getUserDto.PasswordHash,
		CreatedAt:    timestamppb.New(getUserDto.CreatedAt),
	}

	if getUserDto.UpdateAt != nil {
		dist.UpdatedAt = timestamppb.New(*getUserDto.UpdateAt)
	}

	return dist
}

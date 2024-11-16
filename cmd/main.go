package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"user/internal/converters"
	"user/internal/dto"
	userRepository "user/internal/repository/user"
	service "user/internal/services"
	requestsValidation "user/internal/validation/requests"
	desc "user/pkg/user_v1"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	GRPC_PORT = 50051
)

type server struct {
	desc.UnimplementedUserServiceV1Server
	service service.IUserService
}

func (s *server) CreateUser(context context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	validateReq := &requestsValidation.CreateUserRequestWrapper{
		Firstname: req.User.GetFirstname(),
		Lastname:  req.User.GetLastname(),
		Email:     req.User.GetEmail(),
		Password:  req.User.GetPassword(),
	}
	if err := validateReq.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	userId, err := s.service.Create(context, converters.CreateUserRequestToCreateUserDto(req))
	if err != nil {
		return nil, err
	}

	return &desc.CreateUserResponse{
		Id: userId,
	}, nil
}

func (s *server) UpdateUser(context context.Context, req *desc.UpdateUserRequest) (*empty.Empty, error) {
	userId := req.Id
	updateUserDto := &dto.UpdateUserDto{
		FirstName: req.Firstname,
		LastName:  req.Lastname,
	}

	err := s.service.Update(context, userId, updateUserDto)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) GetUser(context context.Context, req *desc.GetUserRequest) (*desc.User, error) {
	user, err := s.service.GetOne(context, req.GetId())
	if err != nil {
		return nil, err
	}

	return converters.GetUserDtoToUserRequest(user), nil
}

func (s *server) GetAllUsers(context context.Context, req *desc.GetAllUsersRequest) (*desc.GetAllUsersResponse, error) {
	users, err := s.service.GetAll(context)
	if err != nil {
		return nil, err
	}

	response := &desc.GetAllUsersResponse{
		Users: []*desc.User{},
	}

	for _, user := range users {
		response.Users = append(response.Users, converters.GetUserDtoToUserRequest(user))
	}

	return response, nil
}

func (s *server) DeleteUser(context context.Context, req *desc.DeleteUserRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func NewServer() {}

func main() {
	lis, err := net.Listen("tcp", net.JoinHostPort("localhost", "50051"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	userRepo := &userRepository.UserRepository{}

	desc.RegisterUserServiceV1Server(grpcServer, &server{
		service: service.NewUserService(userRepo),
	})

	fmt.Println("Server has been started")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

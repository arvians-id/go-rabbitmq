package services

import (
	"context"

	"github.com/arvians-id/go-rabbitmq/gateway/api/auth/client"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
)

type UserServiceContract interface {
	ValidateLogin(ctx context.Context, in *pb.GetValidateLoginRequest) (*pb.GetUserResponse, error)
	Create(ctx context.Context, in *pb.CreateUserRequest) (*pb.GetUserResponse, error)
}

type userService struct {
	UserClient client.UserClient
}

func NewUserService(userClient client.UserClient) UserServiceContract {
	return &userService{
		UserClient: userClient,
	}
}

func (service *userService) ValidateLogin(ctx context.Context, in *pb.GetValidateLoginRequest) (*pb.GetUserResponse, error) {
	return service.UserClient.Client.ValidateLogin(ctx, in)
}

func (service *userService) Create(ctx context.Context, in *pb.CreateUserRequest) (*pb.GetUserResponse, error) {
	return service.UserClient.Client.Create(ctx, in)
}

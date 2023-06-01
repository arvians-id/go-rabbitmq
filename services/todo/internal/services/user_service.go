package services

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/todo/internal/client"
	"github.com/arvians-id/go-rabbitmq/todo/pb"
)

type UserServiceContract interface {
	FindByID(ctx context.Context, in *pb.GetUserByIDRequest) (*pb.GetUserResponse, error)
}

type userService struct {
	UserClient client.UserClient
}

func NewUserService(userClient client.UserClient) UserServiceContract {
	return &userService{
		UserClient: userClient,
	}
}

func (service *userService) FindByID(ctx context.Context, in *pb.GetUserByIDRequest) (*pb.GetUserResponse, error) {
	return service.UserClient.Client.FindByID(ctx, in)
}

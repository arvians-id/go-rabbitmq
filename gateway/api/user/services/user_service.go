package services

import (
	"context"

	"github.com/arvians-id/go-rabbitmq/gateway/api/user/client"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceContract interface {
	FindAll(ctx context.Context, in *emptypb.Empty) (*pb.ListUserResponse, error)
	FindByID(ctx context.Context, in *pb.GetUserByIDRequest) (*pb.GetUserResponse, error)
	Create(ctx context.Context, in *pb.CreateUserRequest) (*pb.GetUserResponse, error)
	Update(ctx context.Context, in *pb.UpdateUserRequest) (*pb.GetUserResponse, error)
	Delete(ctx context.Context, in *pb.GetUserByIDRequest) (*emptypb.Empty, error)
}

type userService struct {
	UserClient client.UserClient
}

func NewUserService(userClient client.UserClient) UserServiceContract {
	return &userService{
		UserClient: userClient,
	}
}

func (service *userService) FindAll(ctx context.Context, in *emptypb.Empty) (*pb.ListUserResponse, error) {
	return service.UserClient.Client.FindAll(ctx, in)
}

func (service *userService) FindByID(ctx context.Context, in *pb.GetUserByIDRequest) (*pb.GetUserResponse, error) {
	return service.UserClient.Client.FindByID(ctx, in)
}

func (service *userService) Create(ctx context.Context, in *pb.CreateUserRequest) (*pb.GetUserResponse, error) {
	return service.UserClient.Client.Create(ctx, in)
}

func (service *userService) Update(ctx context.Context, in *pb.UpdateUserRequest) (*pb.GetUserResponse, error) {
	return service.UserClient.Client.Update(ctx, in)
}

func (service *userService) Delete(ctx context.Context, in *pb.GetUserByIDRequest) (*emptypb.Empty, error) {
	return service.UserClient.Client.Delete(ctx, in)
}

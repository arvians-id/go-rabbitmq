package services

import (
	"context"

	"github.com/arvians-id/go-rabbitmq/gateway/api/user/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/user/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService interface {
	FindAll(ctx context.Context, in *emptypb.Empty) *pb.ListUserResponse
	FindById(ctx context.Context, in *pb.GetUserByIDRequest) *pb.GetUserResponse
	Create(ctx context.Context, in *pb.CreateUserRequest) *pb.GetUserResponse
	Update(ctx context.Context, in *pb.UpdateUserRequest) *pb.GetUserResponse
	Delete(ctx context.Context, in *pb.GetUserByIDRequest) *emptypb.Empty
}

type userService struct {
	UserClient client.UserClient
}

func NewUserService(userClient client.UserClient) UserService {
	return &userService{
		UserClient: userClient,
	}
}

func (userService *userService) FindAll(ctx context.Context, in *emptypb.Empty) *pb.ListUserResponse {
	return nil
}

func (userService *userService) FindById(ctx context.Context, in *pb.GetUserByIDRequest) *pb.GetUserResponse {
	return nil
}

func (userService *userService) Create(ctx context.Context, in *pb.CreateUserRequest) *pb.GetUserResponse {
	return nil
}

func (userService *userService) Update(ctx context.Context, in *pb.UpdateUserRequest) *pb.GetUserResponse {
	return nil
}

func (userService *userService) Delete(ctx context.Context, in *pb.GetUserByIDRequest) *emptypb.Empty {
	return nil
}

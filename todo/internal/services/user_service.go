package services

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/todo/internal/client"
	"github.com/arvians-id/go-rabbitmq/todo/pb"
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
	users, err := service.UserClient.Client.FindAll(ctx, in)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (service *userService) FindByID(ctx context.Context, in *pb.GetUserByIDRequest) (*pb.GetUserResponse, error) {
	user, err := service.UserClient.Client.FindByID(ctx, in)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *userService) Create(ctx context.Context, in *pb.CreateUserRequest) (*pb.GetUserResponse, error) {
	user, err := service.UserClient.Client.Create(ctx, in)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *userService) Update(ctx context.Context, in *pb.UpdateUserRequest) (*pb.GetUserResponse, error) {
	user, err := service.UserClient.Client.Update(ctx, in)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *userService) Delete(ctx context.Context, in *pb.GetUserByIDRequest) (*emptypb.Empty, error) {
	_, err := service.UserClient.Client.Delete(ctx, in)
	if err != nil {
		return nil, err
	}

	return new(emptypb.Empty), nil
}

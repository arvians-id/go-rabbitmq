package services

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/client"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TodoServiceContract interface {
	FindAll(ctx context.Context, in *emptypb.Empty) (*pb.ListTodoResponse, error)
	FindByID(ctx context.Context, in *pb.GetTodoByIDRequest) (*pb.GetTodoResponse, error)
	Create(ctx context.Context, in *pb.CreateTodoRequest) (*pb.GetTodoResponse, error)
	Update(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.GetTodoResponse, error)
	Delete(ctx context.Context, in *pb.GetTodoByIDRequest) (*emptypb.Empty, error)
}

type todoService struct {
	TodoClient client.TodoClient
}

func NewTodoService(todoClient client.TodoClient) TodoServiceContract {
	return &todoService{
		TodoClient: todoClient,
	}
}

func (service *todoService) FindAll(ctx context.Context, in *emptypb.Empty) (*pb.ListTodoResponse, error) {
	return service.TodoClient.Client.FindAll(ctx, in)
}

func (service *todoService) FindByID(ctx context.Context, in *pb.GetTodoByIDRequest) (*pb.GetTodoResponse, error) {
	return service.TodoClient.Client.FindByID(ctx, in)
}

func (service *todoService) Create(ctx context.Context, in *pb.CreateTodoRequest) (*pb.GetTodoResponse, error) {
	return service.TodoClient.Client.Create(ctx, in)
}

func (service *todoService) Update(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.GetTodoResponse, error) {
	return service.TodoClient.Client.Update(ctx, in)
}

func (service *todoService) Delete(ctx context.Context, in *pb.GetTodoByIDRequest) (*emptypb.Empty, error) {
	return service.TodoClient.Client.Delete(ctx, in)
}

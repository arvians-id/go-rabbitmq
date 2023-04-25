package services

import (
	"context"

	"github.com/arvians-id/go-rabbitmq/todo/internal/client"
	"github.com/arvians-id/go-rabbitmq/todo/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryTodoService interface {
	FindAll(ctx context.Context, in *emptypb.Empty) (*pb.ListCategoryTodoResponse, error)
	FindByID(ctx context.Context, in *pb.GetCategoryTodoByIDRequest) (*pb.GetCategoryTodoResponse, error)
	Create(ctx context.Context, in *pb.CreateCategoryTodoRequest) (*pb.GetCategoryTodoResponse, error)
	Delete(ctx context.Context, in *pb.GetCategoryTodoByIDRequest) (*emptypb.Empty, error)
}

type categoryTodoService struct {
	CategoryTodoClient client.CategoryTodoClient
}

func NewCategoryTodoService(categoryTodoClient client.CategoryTodoClient) CategoryTodoService {
	return &categoryTodoService{
		CategoryTodoClient: categoryTodoClient,
	}
}

func (service *categoryTodoService) FindAll(ctx context.Context, in *emptypb.Empty) (*pb.ListCategoryTodoResponse, error) {
	categoryTodos, err := service.CategoryTodoClient.Client.FindAll(ctx, in)
	if err != nil {
		return nil, err
	}

	return categoryTodos, nil
}

func (service *categoryTodoService) FindByID(ctx context.Context, in *pb.GetCategoryTodoByIDRequest) (*pb.GetCategoryTodoResponse, error) {
	categoryTodo, err := service.CategoryTodoClient.Client.FindByID(ctx, in)
	if err != nil {
		return nil, err
	}

	return categoryTodo, err
}

func (service *categoryTodoService) Create(ctx context.Context, in *pb.CreateCategoryTodoRequest) (*pb.GetCategoryTodoResponse, error) {
	categoryTodo, err := service.CategoryTodoClient.Client.Create(ctx, in)
	if err != nil {
		return nil, err
	}

	return categoryTodo, nil
}

func (service *categoryTodoService) Delete(ctx context.Context, in *pb.GetCategoryTodoByIDRequest) (*emptypb.Empty, error) {
	_, err := service.CategoryTodoClient.Client.Delete(ctx, in)
	if err != nil {
		return nil, err
	}

	return new(emptypb.Empty), nil
}
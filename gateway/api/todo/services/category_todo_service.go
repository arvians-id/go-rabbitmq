package services

import (
	"context"

	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryTodoServiceContract interface {
	FindAll(ctx context.Context, in *emptypb.Empty) (*pb.ListCategoryTodoResponse, error)
	FindByID(ctx context.Context, in *pb.GetCategoryTodoByIDRequest) (*pb.GetCategoryTodoResponse, error)
	Create(ctx context.Context, in *pb.CreateCategoryTodoRequest) (*pb.GetCategoryTodoResponse, error)
	Delete(ctx context.Context, in *pb.GetCategoryTodoByIDRequest) (*emptypb.Empty, error)
}

type categoryTodoService struct {
	CategoryTodoClient client.CategoryTodoClient
}

func NewCategoryTodoService(categoryTodoClient client.CategoryTodoClient) CategoryTodoServiceContract {
	return &categoryTodoService{
		CategoryTodoClient: categoryTodoClient,
	}
}

func (service *categoryTodoService) FindAll(ctx context.Context, in *emptypb.Empty) (*pb.ListCategoryTodoResponse, error) {
	return service.CategoryTodoClient.Client.FindAll(ctx, in)
}

func (service *categoryTodoService) FindByID(ctx context.Context, in *pb.GetCategoryTodoByIDRequest) (*pb.GetCategoryTodoResponse, error) {
	return service.CategoryTodoClient.Client.FindByID(ctx, in)
}

func (service *categoryTodoService) Create(ctx context.Context, in *pb.CreateCategoryTodoRequest) (*pb.GetCategoryTodoResponse, error) {
	return service.CategoryTodoClient.Client.Create(ctx, in)
}

func (service *categoryTodoService) Delete(ctx context.Context, in *pb.GetCategoryTodoByIDRequest) (*emptypb.Empty, error) {
	return service.CategoryTodoClient.Client.Delete(ctx, in)
}

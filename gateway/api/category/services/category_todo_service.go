package services

import (
	"context"

	"github.com/arvians-id/go-rabbitmq/gateway/api/category/client"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryServiceContract interface {
	FindAll(ctx context.Context, in *emptypb.Empty) (*pb.ListCategoryResponse, error)
	FindByID(ctx context.Context, in *pb.GetCategoryByIDRequest) (*pb.GetCategoryResponse, error)
	Create(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.GetCategoryResponse, error)
	Delete(ctx context.Context, in *pb.GetCategoryByIDRequest) (*emptypb.Empty, error)
}

type categoryService struct {
	CategoryClient client.CategoryClient
}

func NewCategoryService(categoryClient client.CategoryClient) CategoryServiceContract {
	return &categoryService{
		CategoryClient: categoryClient,
	}
}

func (service *categoryService) FindAll(ctx context.Context, in *emptypb.Empty) (*pb.ListCategoryResponse, error) {
	return service.CategoryClient.Client.FindAll(ctx, in)
}

func (service *categoryService) FindByID(ctx context.Context, in *pb.GetCategoryByIDRequest) (*pb.GetCategoryResponse, error) {
	return service.CategoryClient.Client.FindByID(ctx, in)
}

func (service *categoryService) Create(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.GetCategoryResponse, error) {
	return service.CategoryClient.Client.Create(ctx, in)
}

func (service *categoryService) Delete(ctx context.Context, in *pb.GetCategoryByIDRequest) (*emptypb.Empty, error) {
	return service.CategoryClient.Client.Delete(ctx, in)
}

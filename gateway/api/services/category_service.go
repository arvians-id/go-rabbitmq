package services

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/client"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryServiceContract interface {
	FindAll(ctx context.Context) (*pb.ListCategoryResponse, int, error)
	FindByTodoIDs(ctx context.Context, in *pb.GetCategoryByTodoIDsRequest) (*pb.ListCategoryWithTodoIDResponse, int, error)
	FindByID(ctx context.Context, in *pb.GetCategoryByIDRequest) (*pb.GetCategoryResponse, int, error)
	Create(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.GetCategoryResponse, int, error)
	Delete(ctx context.Context, in *pb.GetCategoryByIDRequest) (int, error)
}

type categoryService struct {
	CategoryClient client.CategoryClient
}

func NewCategoryService(categoryClient client.CategoryClient) CategoryServiceContract {
	return &categoryService{
		CategoryClient: categoryClient,
	}
}

func (service *categoryService) FindAll(ctx context.Context) (*pb.ListCategoryResponse, int, error) {
	categories, err := service.CategoryClient.Client.FindAll(ctx, new(emptypb.Empty))
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return categories, fiber.StatusOK, nil
}

func (service *categoryService) FindByTodoIDs(ctx context.Context, in *pb.GetCategoryByTodoIDsRequest) (*pb.ListCategoryWithTodoIDResponse, int, error) {
	categories, err := service.CategoryClient.Client.FindByTodoIDs(ctx, in)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return categories, fiber.StatusOK, nil
}

func (service *categoryService) FindByID(ctx context.Context, in *pb.GetCategoryByIDRequest) (*pb.GetCategoryResponse, int, error) {
	category, err := service.CategoryClient.Client.FindByID(ctx, in)
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return nil, fiber.StatusNotFound, err
		}
		return nil, fiber.StatusInternalServerError, err
	}

	return category, fiber.StatusOK, nil
}

func (service *categoryService) Create(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.GetCategoryResponse, int, error) {
	category, err := service.CategoryClient.Client.Create(ctx, in)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return category, fiber.StatusCreated, nil
}

func (service *categoryService) Delete(ctx context.Context, in *pb.GetCategoryByIDRequest) (int, error) {
	_, err := service.CategoryClient.Client.Delete(ctx, in)
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return fiber.StatusNotFound, err
		}
		return fiber.StatusInternalServerError, err
	}

	return fiber.StatusOK, nil
}

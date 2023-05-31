package services

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/todo/dto"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

type TodoServiceContract interface {
	DisplayTodoCategoryList(ctx context.Context) (*dto.DisplayCategoryTodoListResponse, int, error)
	FindAll(ctx context.Context) (*pb.ListTodoResponse, int, error)
	FindByUserIDs(ctx context.Context, in *pb.GetTodoByUserIDsRequest) (*pb.ListTodoResponse, int, error)
	FindByID(ctx context.Context, in *pb.GetTodoByIDRequest) (*pb.GetTodoResponse, int, error)
	Create(ctx context.Context, in *pb.CreateTodoRequest) (*pb.GetTodoResponse, int, error)
	Update(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.GetTodoResponse, int, error)
	Delete(ctx context.Context, in *pb.GetTodoByIDRequest) (int, error)
}

type todoService struct {
	TodoClient     client.TodoClient
	CategoryClient client.CategoryClient
}

func NewTodoService(todoClient client.TodoClient, categoryClient client.CategoryClient) TodoServiceContract {
	return &todoService{
		TodoClient:     todoClient,
		CategoryClient: categoryClient,
	}
}

func (service *todoService) DisplayTodoCategoryList(ctx context.Context) (*dto.DisplayCategoryTodoListResponse, int, error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var todos *pb.ListTodoResponse
	var categories *pb.ListCategoryResponse
	var err error
	wg.Add(2)

	go func() {
		var errGo error
		todos, errGo = service.TodoClient.Client.FindAll(ctx, new(emptypb.Empty))
		if errGo != nil {
			mutex.Lock()
			err = errGo
			mutex.Unlock()
		}
		defer wg.Done()
	}()

	go func() {
		var errGo error
		categories, errGo = service.CategoryClient.Client.FindAll(ctx, new(emptypb.Empty))
		if errGo != nil {
			mutex.Lock()
			err = errGo
			mutex.Unlock()
		}
		defer wg.Done()
	}()
	wg.Wait()

	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return &dto.DisplayCategoryTodoListResponse{
		Todos:      todos.GetTodos(),
		Categories: categories.GetCategories(),
	}, fiber.StatusOK, nil
}

func (service *todoService) FindAll(ctx context.Context) (*pb.ListTodoResponse, int, error) {
	todos, err := service.TodoClient.Client.FindAll(ctx, new(emptypb.Empty))
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return todos, fiber.StatusOK, nil
}

func (service *todoService) FindByUserIDs(ctx context.Context, in *pb.GetTodoByUserIDsRequest) (*pb.ListTodoResponse, int, error) {
	todos, err := service.TodoClient.Client.FindByUserIDs(ctx, in)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return todos, fiber.StatusOK, nil
}

func (service *todoService) FindByID(ctx context.Context, in *pb.GetTodoByIDRequest) (*pb.GetTodoResponse, int, error) {
	todo, err := service.TodoClient.Client.FindByID(ctx, in)
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return nil, fiber.StatusNotFound, err
		}
		return nil, fiber.StatusInternalServerError, err
	}

	return todo, fiber.StatusOK, nil
}

func (service *todoService) Create(ctx context.Context, in *pb.CreateTodoRequest) (*pb.GetTodoResponse, int, error) {
	todo, err := service.TodoClient.Client.Create(ctx, in)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return todo, fiber.StatusCreated, nil
}

func (service *todoService) Update(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.GetTodoResponse, int, error) {
	todo, err := service.TodoClient.Client.Update(ctx, in)
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return nil, fiber.StatusNotFound, err
		}
		return nil, fiber.StatusInternalServerError, err
	}

	return todo, fiber.StatusOK, nil
}

func (service *todoService) Delete(ctx context.Context, in *pb.GetTodoByIDRequest) (int, error) {
	_, err := service.TodoClient.Client.Delete(ctx, in)
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return fiber.StatusNotFound, err
		}
		return fiber.StatusInternalServerError, err
	}

	return fiber.StatusOK, nil
}

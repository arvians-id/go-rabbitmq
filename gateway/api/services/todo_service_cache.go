package services

import (
	"context"
	"fmt"
	"github.com/arvians-id/go-rabbitmq/gateway/api/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/todo/dto"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
	"time"
)

type todoServiceCache struct {
	TodoClient     client.TodoClient
	RedisClient    *redis.Client
	CategoryClient CategoryServiceContract
}

func NewTodoServiceCache(todoClient client.TodoClient, redisClient *redis.Client, categoryService CategoryServiceContract) TodoServiceContract {
	return &todoServiceCache{
		TodoClient:     todoClient,
		RedisClient:    redisClient,
		CategoryClient: categoryService,
	}
}

func (service *todoServiceCache) DisplayTodoCategoryList(ctx context.Context) (*dto.DisplayCategoryTodoListResponse, int, error) {
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
		categories, _, errGo = service.CategoryClient.FindAll(ctx, new(emptypb.Empty))
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

func (service *todoServiceCache) FindAll(ctx context.Context, in *emptypb.Empty) (*pb.ListTodoResponse, int, error) {
	todosCached, err := service.RedisClient.Get(ctx, "todos").Bytes()
	if err == redis.Nil {
		todos, err := service.TodoClient.Client.FindAll(ctx, in)
		if err != nil {
			return nil, fiber.StatusInternalServerError, err
		}

		todosJSON, err := json.Marshal(todos)
		if err != nil {
			return nil, fiber.StatusBadRequest, err
		}

		err = service.RedisClient.Set(ctx, "todos", todosJSON, time.Second*20).Err()
		if err != nil {
			return nil, fiber.StatusBadRequest, err
		}

		return todos, fiber.StatusOK, nil
	}

	var todos *pb.ListTodoResponse
	err = json.Unmarshal(todosCached, &todos)
	if err != nil {
		return nil, fiber.StatusBadRequest, err
	}

	return todos, fiber.StatusOK, nil
}

func (service *todoServiceCache) FindByID(ctx context.Context, in *pb.GetTodoByIDRequest) (*pb.GetTodoResponse, int, error) {
	keys := fmt.Sprintf("todo:%d", in.GetId())
	todoCached, err := service.RedisClient.Get(ctx, keys).Bytes()
	if err == redis.Nil {
		todo, err := service.TodoClient.Client.FindByID(ctx, in)
		if err != nil {
			if err.Error() == response.GrpcErrorNotFound {
				return nil, fiber.StatusNotFound, err
			}
			return nil, fiber.StatusInternalServerError, err
		}

		todoJSON, err := json.Marshal(todo)
		if err != nil {
			return nil, fiber.StatusBadRequest, err
		}

		err = service.RedisClient.Set(ctx, keys, todoJSON, time.Second*10).Err()
		if err != nil {
			return nil, fiber.StatusBadRequest, err
		}

		return todo, fiber.StatusOK, nil
	}

	var todo *pb.GetTodoResponse
	err = json.Unmarshal(todoCached, &todo)
	if err != nil {
		return nil, fiber.StatusBadRequest, err
	}

	return todo, fiber.StatusOK, nil
}

func (service *todoServiceCache) Create(ctx context.Context, in *pb.CreateTodoRequest) (*pb.GetTodoResponse, int, error) {
	todo, err := service.TodoClient.Client.Create(ctx, in)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	err = service.RedisClient.Del(ctx, "todos").Err()
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return todo, fiber.StatusCreated, nil
}

func (service *todoServiceCache) Update(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.GetTodoResponse, int, error) {
	todo, err := service.TodoClient.Client.Update(ctx, in)
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return nil, fiber.StatusNotFound, err
		}
		return nil, fiber.StatusInternalServerError, err
	}

	keys := fmt.Sprintf("todo:%d", in.GetId())
	err = service.RedisClient.Del(ctx, keys, "todos").Err()
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return todo, fiber.StatusOK, nil
}

func (service *todoServiceCache) Delete(ctx context.Context, in *pb.GetTodoByIDRequest) (*emptypb.Empty, int, error) {
	_, err := service.TodoClient.Client.Delete(ctx, in)
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return nil, fiber.StatusNotFound, err
		}
		return nil, fiber.StatusInternalServerError, err
	}

	keys := fmt.Sprintf("todo:%d", in.GetId())
	err = service.RedisClient.Del(ctx, keys, "todos").Err()
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return new(emptypb.Empty), fiber.StatusOK, err
}

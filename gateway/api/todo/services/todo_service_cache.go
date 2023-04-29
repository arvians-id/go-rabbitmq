package services

import (
	"context"
	"fmt"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/pb"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type todoServiceCache struct {
	TodoClient  client.TodoClient
	RedisClient *redis.Client
}

func NewTodoServiceCache(todoClient client.TodoClient, redisClient *redis.Client) TodoServiceContract {
	return &todoServiceCache{
		TodoClient:  todoClient,
		RedisClient: redisClient,
	}
}

func (service *todoServiceCache) FindAll(ctx context.Context, in *emptypb.Empty) (*pb.ListTodoResponse, error) {
	todosCached, err := service.RedisClient.Get(ctx, "todos").Bytes()
	if err == redis.Nil {
		todos, err := service.TodoClient.Client.FindAll(ctx, in)
		if err != nil {
			return nil, err
		}

		todosJSON, err := json.Marshal(todos)
		if err != nil {
			return nil, err
		}

		err = service.RedisClient.Set(ctx, "todos", todosJSON, time.Second*20).Err()
		if err != nil {
			return nil, err
		}

		return todos, nil
	}

	var todos *pb.ListTodoResponse
	err = json.Unmarshal(todosCached, &todos)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (service *todoServiceCache) FindByID(ctx context.Context, in *pb.GetTodoByIDRequest) (*pb.GetTodoResponse, error) {
	keys := fmt.Sprintf("todo:%d", in.GetId())
	todoCached, err := service.RedisClient.Get(ctx, keys).Bytes()
	if err == redis.Nil {
		todo, err := service.TodoClient.Client.FindByID(ctx, in)
		if err != nil {
			return nil, err
		}

		todoJSON, err := json.Marshal(todo)
		if err != nil {
			return nil, err
		}

		err = service.RedisClient.Set(ctx, keys, todoJSON, time.Second*10).Err()
		if err != nil {
			return nil, err
		}

		return todo, nil
	}

	var todo pb.GetTodoResponse
	err = json.Unmarshal(todoCached, &todo)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (service *todoServiceCache) Create(ctx context.Context, in *pb.CreateTodoRequest) (*pb.GetTodoResponse, error) {
	todo, err := service.TodoClient.Client.Create(ctx, in)
	if err != nil {
		return nil, err
	}

	err = service.RedisClient.Del(ctx, "todos").Err()
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (service *todoServiceCache) Update(ctx context.Context, in *pb.UpdateTodoRequest) (*pb.GetTodoResponse, error) {
	todo, err := service.TodoClient.Client.Update(ctx, in)
	if err != nil {
		return nil, err
	}

	keys := fmt.Sprintf("todo:%d", in.GetId())
	err = service.RedisClient.Del(ctx, keys, "todos").Err()
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (service *todoServiceCache) Delete(ctx context.Context, in *pb.GetTodoByIDRequest) (*emptypb.Empty, error) {
	_, err := service.TodoClient.Client.Delete(ctx, in)
	if err != nil {
		return nil, err
	}

	keys := fmt.Sprintf("todo:%d", in.GetId())
	err = service.RedisClient.Del(ctx, keys, "todos").Err()
	if err != nil {
		return nil, err
	}

	return new(emptypb.Empty), err
}

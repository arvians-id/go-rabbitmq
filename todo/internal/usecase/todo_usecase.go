package usecase

import (
	"context"
	"time"

	"github.com/arvians-id/go-rabbitmq/todo/internal/model"
	"github.com/arvians-id/go-rabbitmq/todo/internal/repository"
	"github.com/arvians-id/go-rabbitmq/todo/internal/services"
	"github.com/arvians-id/go-rabbitmq/todo/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TodoUsecase struct {
	UserService    services.UserService
	TodoRepository repository.TodoRepository
	pb.UnimplementedTodoServiceServer
}

func NewTodoUsecase(userService services.UserService, todoRepository repository.TodoRepository) pb.TodoServiceServer {
	return &TodoUsecase{
		UserService:    userService,
		TodoRepository: todoRepository,
	}
}

func (usecase *TodoUsecase) FindAll(ctx context.Context, empty *emptypb.Empty) (*pb.ListTodoResponse, error) {
	todos, err := usecase.TodoRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var todosPB []*pb.Todo
	for _, todo := range todos {
		todosPB = append(todosPB, todo.ToPB())
	}

	return &pb.ListTodoResponse{
		Todos: todosPB,
	}, nil
}

func (usecase *TodoUsecase) FindByID(ctx context.Context, req *pb.GetTodoByIDRequest) (*pb.GetTodoResponse, error) {
	todo, err := usecase.TodoRepository.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetTodoResponse{
		Todo: todo.ToPB(),
	}, nil
}

func (usecase *TodoUsecase) Create(ctx context.Context, req *pb.CreateTodoRequest) (*pb.GetTodoResponse, error) {
	userCheck, err := usecase.UserService.FindByID(ctx, &pb.GetUserByIDRequest{
		Id: req.GetUserId(),
	})
	if err != nil {
		return nil, err
	}

	todoCreated, err := usecase.TodoRepository.Create(ctx, &model.Todo{
		Title:          req.GetTitle(),
		Description:    req.GetDescription(),
		UserId:         userCheck.User.GetId(),
		CategoryTodoId: req.GetCategoryTodoId(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetTodoResponse{
		Todo: todoCreated.ToPB(),
	}, nil
}

func (usecase *TodoUsecase) Update(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.GetTodoResponse, error) {
	todoCheck, err := usecase.TodoRepository.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	todoUpdated, err := usecase.TodoRepository.Update(ctx, &model.Todo{
		Id:          todoCheck.Id,
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		IsDone:      req.IsDone,
		UserId:      req.GetUserId(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetTodoResponse{
		Todo: todoUpdated.ToPB(),
	}, nil
}

func (usecase *TodoUsecase) Delete(ctx context.Context, req *pb.GetTodoByIDRequest) (*emptypb.Empty, error) {
	todoCheck, err := usecase.TodoRepository.FindByID(ctx, req.Id)
	if err != nil {
		return new(emptypb.Empty), nil
	}

	err = usecase.TodoRepository.Delete(ctx, todoCheck.Id)
	if err != nil {
		return new(emptypb.Empty), nil
	}

	return new(emptypb.Empty), nil
}

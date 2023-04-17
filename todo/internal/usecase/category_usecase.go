package usecase

import (
	"context"
	"time"

	"github.com/arvians-id/go-rabbitmq/todo/internal/model"
	"github.com/arvians-id/go-rabbitmq/todo/internal/repository"
	"github.com/arvians-id/go-rabbitmq/todo/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryTodoUsecase struct {
	CategoryTodoRepository repository.CategoryTodoRepository
	pb.UnimplementedCategoryTodoServiceServer
}

func NewCategoryTodoUsecase(todoRepository repository.CategoryTodoRepository) pb.CategoryTodoServiceServer {
	return &CategoryTodoUsecase{
		CategoryTodoRepository: todoRepository,
	}
}

func (usecase *CategoryTodoUsecase) FindAll(ctx context.Context, empty *emptypb.Empty) (*pb.ListCategoryTodoResponse, error) {
	todos, err := usecase.CategoryTodoRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var todosPB []*pb.CategoryTodo
	for _, todo := range todos {
		todosPB = append(todosPB, todo.ToPB())
	}

	return &pb.ListCategoryTodoResponse{
		CategoryTodos: todosPB,
	}, nil
}

func (usecase *CategoryTodoUsecase) FindByID(ctx context.Context, req *pb.GetCategoryTodoByIDRequest) (*pb.GetCategoryTodoResponse, error) {
	todo, err := usecase.CategoryTodoRepository.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetCategoryTodoResponse{
		CategoryTodo: todo.ToPB(),
	}, nil
}

func (usecase *CategoryTodoUsecase) Create(ctx context.Context, req *pb.CreateCategoryTodoRequest) (*pb.GetCategoryTodoResponse, error) {
	todoCreated, err := usecase.CategoryTodoRepository.Create(ctx, &model.CategoryTodo{
		Name:      req.GetName(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetCategoryTodoResponse{
		CategoryTodo: todoCreated.ToPB(),
	}, nil
}

func (usecase *CategoryTodoUsecase) Delete(ctx context.Context, req *pb.GetCategoryTodoByIDRequest) (*emptypb.Empty, error) {
	todoCheck, err := usecase.CategoryTodoRepository.FindByID(ctx, req.Id)
	if err != nil {
		return new(emptypb.Empty), nil
	}

	err = usecase.CategoryTodoRepository.Delete(ctx, todoCheck.Id)
	if err != nil {
		return new(emptypb.Empty), nil
	}

	return new(emptypb.Empty), nil
}

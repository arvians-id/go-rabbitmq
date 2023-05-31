package usecase

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/category/internal/model"
	"github.com/arvians-id/go-rabbitmq/category/internal/repository"
	"github.com/arvians-id/go-rabbitmq/category/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryUsecase struct {
	CategoryRepository repository.CategoryRepository
	pb.UnimplementedCategoryServiceServer
}

func NewCategoryUsecase(todoRepository repository.CategoryRepository) pb.CategoryServiceServer {
	return &CategoryUsecase{
		CategoryRepository: todoRepository,
	}
}

func (usecase *CategoryUsecase) FindAll(ctx context.Context, empty *emptypb.Empty) (*pb.ListCategoryResponse, error) {
	todos, err := usecase.CategoryRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var todosPB []*pb.Category
	for _, todo := range todos {
		todosPB = append(todosPB, todo.ToPB())
	}

	return &pb.ListCategoryResponse{
		Categories: todosPB,
	}, nil
}

func (usecase *CategoryUsecase) FindByTodoIDs(ctx context.Context, req *pb.GetCategoryByTodoIDsRequest) (*pb.ListCategoryWithTodoIDResponse, error) {
	todos, err := usecase.CategoryRepository.FindByTodoIDs(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	var todosPB []*pb.CategoryWithTodoID
	for _, todo := range todos {
		todosPB = append(todosPB, todo.ToPB())
	}

	return &pb.ListCategoryWithTodoIDResponse{
		Categories: todosPB,
	}, nil
}

func (usecase *CategoryUsecase) FindByID(ctx context.Context, req *pb.GetCategoryByIDRequest) (*pb.GetCategoryResponse, error) {
	todo, err := usecase.CategoryRepository.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetCategoryResponse{
		Category: todo.ToPB(),
	}, nil
}

func (usecase *CategoryUsecase) Create(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.GetCategoryResponse, error) {
	todoCreated, err := usecase.CategoryRepository.Create(ctx, &model.Category{
		Name: req.GetName(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetCategoryResponse{
		Category: todoCreated.ToPB(),
	}, nil
}

func (usecase *CategoryUsecase) Delete(ctx context.Context, req *pb.GetCategoryByIDRequest) (*emptypb.Empty, error) {
	todoCheck, err := usecase.CategoryRepository.FindByID(ctx, req.Id)
	if err != nil {
		return new(emptypb.Empty), err
	}

	err = usecase.CategoryRepository.Delete(ctx, todoCheck.Id)
	if err != nil {
		return new(emptypb.Empty), err
	}

	return new(emptypb.Empty), nil
}

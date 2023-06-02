package resolver

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/gql/model"
	"github.com/arvians-id/go-rabbitmq/gateway/api/middleware"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
)

func (r *queryResolver) TodoFindAll(ctx context.Context) ([]*model.Todo, error) {
	todos, _, err := r.TodoService.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.Todo
	for _, todo := range todos.Todos {
		result = append(result, &model.Todo{
			Id:          todo.GetId(),
			Title:       todo.GetTitle(),
			Description: todo.GetDescription(),
			IsDone:      todo.GetIsDone(),
			UserId:      todo.GetUserId(),
			CreatedAt:   todo.GetCreatedAt(),
			UpdatedAt:   todo.GetUpdatedAt(),
		})
	}

	return result, nil
}

func (r *queryResolver) TodoFindByID(ctx context.Context, id int64) (*model.Todo, error) {
	todo, _, err := r.TodoService.FindByID(ctx, &pb.GetTodoByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &model.Todo{
		Id:          todo.GetTodo().GetId(),
		Title:       todo.GetTodo().GetTitle(),
		Description: todo.GetTodo().GetDescription(),
		IsDone:      todo.GetTodo().GetIsDone(),
		UserId:      todo.GetTodo().GetUserId(),
		CreatedAt:   todo.GetTodo().GetCreatedAt(),
		UpdatedAt:   todo.GetTodo().GetUpdatedAt(),
	}, nil
}

func (r *mutationResolver) TodoCreate(ctx context.Context, input model.TodoCreateRequest) (*model.Todo, error) {
	err := helper.ValidateStruct(input)
	if err != nil {
		return nil, err
	}

	todo, _, err := r.TodoService.Create(ctx, &pb.CreateTodoRequest{
		Title:       input.Title,
		Description: input.Description,
		UserId:      input.UserId,
		CategoryId:  input.Categories,
	})
	if err != nil {
		return nil, err
	}

	return &model.Todo{
		Id:          todo.GetTodo().GetId(),
		Title:       todo.GetTodo().GetTitle(),
		Description: todo.GetTodo().GetDescription(),
		IsDone:      todo.GetTodo().GetIsDone(),
		UserId:      todo.GetTodo().GetUserId(),
		CreatedAt:   todo.GetTodo().GetCreatedAt(),
		UpdatedAt:   todo.GetTodo().GetUpdatedAt(),
	}, nil
}

func (r *mutationResolver) TodoUpdate(ctx context.Context, id int64, input model.TodoUpdateRequest) (*model.Todo, error) {
	err := helper.ValidateStruct(input)
	if err != nil {
		return nil, err
	}

	todo, _, err := r.TodoService.Update(ctx, &pb.UpdateTodoRequest{
		Id:          id,
		Title:       input.Title,
		Description: input.Description,
		IsDone:      &input.IsDone,
		UserId:      input.UserId,
		CategoryId:  input.Categories,
	})
	if err != nil {
		return nil, err
	}

	return &model.Todo{
		Id:          todo.GetTodo().GetId(),
		Title:       todo.GetTodo().GetTitle(),
		Description: todo.GetTodo().GetDescription(),
		IsDone:      todo.GetTodo().GetIsDone(),
		UserId:      todo.GetTodo().GetUserId(),
		CreatedAt:   todo.GetTodo().GetCreatedAt(),
		UpdatedAt:   todo.GetTodo().GetUpdatedAt(),
	}, nil
}

func (r *mutationResolver) TodoDelete(ctx context.Context, id int64) (bool, error) {
	_, err := r.TodoService.Delete(ctx, &pb.GetTodoByIDRequest{
		Id: id,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (t *todoResolver) Categories(ctx context.Context, obj *model.Todo) ([]*model.Category, error) {
	categories, err := middleware.GetLoaders(ctx).CategoryServiceFindByTodoIDs.Load(obj.Id)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (t *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	users, err := middleware.GetLoaders(ctx).UserServiceFindByIDs.Load(obj.UserId)
	if err != nil {
		return nil, err
	}

	return users, nil
}

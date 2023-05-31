package resolver

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/gql/model"
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
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) TodoCreate(ctx context.Context, input model.TodoCreateRequest) (*model.Todo, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) TodoUpdate(ctx context.Context, id int64, input model.TodoUpdateRequest) (*model.Todo, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) TodoDelete(ctx context.Context, id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (t *todoResolver) Categories(ctx context.Context, obj *model.Todo) ([]*model.Category, error) {
	categories, _, err := t.CategoryServices.FindAllByTodoID(ctx, &pb.GetCategoryByTodoIDRequest{
		Id: obj.Id,
	})
	if err != nil {
		return nil, err
	}

	var result []*model.Category
	for _, category := range categories.Categories {
		result = append(result, &model.Category{
			Id:        category.GetId(),
			Name:      category.GetName(),
			CreatedAt: category.GetCreatedAt(),
			UpdatedAt: category.GetUpdatedAt(),
		})
	}

	return result, nil
}

func (t *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	load, err := model.GetUserLoader(ctx).Load(obj.UserId)
	if err != nil {
		return nil, err
	}

	return load, nil
}

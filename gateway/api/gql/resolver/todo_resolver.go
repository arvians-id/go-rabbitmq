package resolver

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/gql/model"
)

func (r *queryResolver) TodoFindAll(ctx context.Context) ([]*model.Todo, error) {
	todos, _, err := r.TodoService.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.Todo
	for _, todo := range todos.Todos {
		var categories []*model.Category
		for _, category := range todo.GetCategories() {
			categories = append(categories, &model.Category{
				Id:        category.GetId(),
				Name:      category.GetName(),
				CreatedAt: category.GetCreatedAt(),
				UpdatedAt: category.GetUpdatedAt(),
			})
		}
		result = append(result, &model.Todo{
			Id:          todo.GetId(),
			Title:       todo.GetTitle(),
			Description: todo.GetDescription(),
			IsDone:      todo.GetIsDone(),
			UserId:      todo.GetUserId(),
			Categories:  categories,
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

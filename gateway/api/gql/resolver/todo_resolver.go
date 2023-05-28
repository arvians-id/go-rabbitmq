package resolver

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/gql/model"
)

func (r *queryResolver) TodoDisplayTodoCategoryList(ctx context.Context) ([]*model.DisplayCategoryTodoListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *queryResolver) TodoFindAll(ctx context.Context) ([]*model.Todo, error) {
	//TODO implement me
	panic("implement me")
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

func (r *mutationResolver) TodoDelete(ctx context.Context, id int64) (*model.Todo, error) {
	//TODO implement me
	panic("implement me")
}

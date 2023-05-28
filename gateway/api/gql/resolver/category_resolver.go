package resolver

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/gql/model"
)

// CategoryFindAll is the resolver for the CategoryFindAll field.
func (r *queryResolver) CategoryFindAll(ctx context.Context) ([]*model.Category, error) {
	panic("not implemented")
}

// CategoryFindByID is the resolver for the CategoryFindById field.
func (r *queryResolver) CategoryFindByID(ctx context.Context, id int64) (*model.Category, error) {
	panic("not implemented")
}

// CategoryCreate is the resolver for the CategoryCreate field.
func (r *mutationResolver) CategoryCreate(ctx context.Context, input model.CategoryCreateRequest) (*model.Todo, error) {
	panic("not implemented")
}

// CategoryDelete is the resolver for the CategoryDelete field.
func (r *mutationResolver) CategoryDelete(ctx context.Context, id int64) (*model.Todo, error) {
	panic("not implemented")
}

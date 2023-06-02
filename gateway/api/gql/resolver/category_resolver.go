package resolver

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/gql/model"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
)

// CategoryFindAll is the resolver for the CategoryFindAll field.
func (r *queryResolver) CategoryFindAll(ctx context.Context) ([]*model.Category, error) {
	categories, _, err := r.CategoryServices.FindAll(ctx)
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

// CategoryFindByID is the resolver for the CategoryFindById field.
func (r *queryResolver) CategoryFindByID(ctx context.Context, id int64) (*model.Category, error) {
	category, _, err := r.CategoryServices.FindByID(ctx, &pb.GetCategoryByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &model.Category{
		Id:        category.GetCategory().GetId(),
		Name:      category.GetCategory().GetName(),
		CreatedAt: category.GetCategory().GetCreatedAt(),
		UpdatedAt: category.GetCategory().GetUpdatedAt(),
	}, nil
}

// CategoryCreate is the resolver for the CategoryCreate field.
func (r *mutationResolver) CategoryCreate(ctx context.Context, input model.CategoryCreateRequest) (*model.Category, error) {
	err := helper.ValidateStruct(input)
	if err != nil {
		return nil, err
	}

	category, _, err := r.CategoryServices.Create(ctx, &pb.CreateCategoryRequest{
		Name: input.Name,
	})
	if err != nil {
		return nil, err
	}

	return &model.Category{
		Id:        category.GetCategory().GetId(),
		Name:      category.GetCategory().GetName(),
		CreatedAt: category.GetCategory().GetCreatedAt(),
		UpdatedAt: category.GetCategory().GetUpdatedAt(),
	}, nil
}

// CategoryDelete is the resolver for the CategoryDelete field.
func (r *mutationResolver) CategoryDelete(ctx context.Context, id int64) (bool, error) {
	_, err := r.CategoryServices.Delete(ctx, &pb.GetCategoryByIDRequest{
		Id: id,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

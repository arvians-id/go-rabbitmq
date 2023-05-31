package resolver

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/gql/model"
	"github.com/arvians-id/go-rabbitmq/gateway/api/middleware"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
)

func (r *queryResolver) UserFindAll(ctx context.Context) ([]*model.User, error) {
	users, _, err := r.UserService.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.User
	for _, user := range users.Users {
		result = append(result, &model.User{
			Id:        user.GetId(),
			Name:      user.GetName(),
			Email:     user.GetEmail(),
			CreatedAt: user.GetCreatedAt(),
			UpdatedAt: user.GetUpdatedAt(),
		})
	}

	return result, nil
}

func (r *queryResolver) UserFindByID(ctx context.Context, id int64) (*model.User, error) {
	user, _, err := r.UserService.FindByID(ctx, &pb.GetUserByIDRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return &model.User{
		Id:        user.User.Id,
		Name:      user.User.Name,
		Email:     user.User.Email,
		CreatedAt: user.User.CreatedAt,
		UpdatedAt: user.User.UpdatedAt,
	}, nil
}

func (r *mutationResolver) UserCreate(ctx context.Context, input model.UserCreateRequest) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) UserUpdate(ctx context.Context, id int64, input model.UserUpdateRequest) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) UserDelete(ctx context.Context, id int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userResolver) Todos(ctx context.Context, obj *model.User) ([]*model.Todo, error) {
	todos, err := middleware.GetLoaders(ctx).TodoServiceFindByUserIDs.Load(obj.Id)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

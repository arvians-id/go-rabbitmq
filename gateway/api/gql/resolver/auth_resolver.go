package resolver

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/gql/model"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
)

func (r *mutationResolver) AuthLogin(ctx context.Context, input model.AuthLoginRequest) (*model.AuthLoginResponse, error) {
	err := helper.ValidateStruct(input)
	if err != nil {
		return nil, err
	}

	signedString, _, err := r.UserService.ValidateLogin(ctx, &pb.GetValidateLoginRequest{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	return &model.AuthLoginResponse{
		Token: signedString,
	}, nil
}

func (r *mutationResolver) AuthRegister(ctx context.Context, input model.AuthRegisterRequest) (*model.AuthRegisterResponse, error) {
	err := helper.ValidateStruct(input)
	if err != nil {
		return nil, err
	}

	user, _, err := r.UserService.Create(ctx, &pb.CreateUserRequest{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	return &model.AuthRegisterResponse{
		Id:        user.GetUser().GetId(),
		Name:      user.GetUser().GetName(),
		Email:     user.GetUser().GetEmail(),
		CreatedAt: user.GetUser().GetCreatedAt(),
		UpdatedAt: user.GetUser().GetUpdatedAt(),
	}, nil
}

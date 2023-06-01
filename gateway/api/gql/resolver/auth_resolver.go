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
	//TODO implement me
	panic("implement me")
}

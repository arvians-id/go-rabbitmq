package resolver

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/gql/model"
)

func (r *mutationResolver) AuthLogin(ctx context.Context, input model.AuthLoginRequest) (*model.AuthLoginResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *mutationResolver) AuthRegister(ctx context.Context, input model.AuthRegisterRequest) (*model.AuthRegisterResponse, error) {
	//TODO implement me
	panic("implement me")
}

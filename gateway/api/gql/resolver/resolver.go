package resolver

import (
	"github.com/99designs/gqlgen/plugin/federation/testdata/entityresolver/generated"
	"github.com/arvians-id/go-rabbitmq/gateway/api/gql"
	"github.com/arvians-id/go-rabbitmq/gateway/api/services"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	UserService      services.UserServiceContract
	TodoService      services.TodoServiceContract
	CategoryServices services.CategoryServiceContract
}

func (r *Resolver) Entity() generated.EntityResolver {
	//TODO implement me
	panic("implement me")
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() gql.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() gql.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

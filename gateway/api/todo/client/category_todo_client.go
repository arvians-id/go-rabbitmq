package client

import (
	"log"

	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"google.golang.org/grpc"
)

type CategoryTodoClient struct {
	Client pb.CategoryTodoServiceClient
}

func InitCategoryTodoClient(configuration config.Config) CategoryTodoClient {
	connection, err := grpc.Dial(configuration.Get("TODO_SERVICE_URL"), grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(config.NewGRPUnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalln(err)
	}

	return CategoryTodoClient{
		Client: pb.NewCategoryTodoServiceClient(connection),
	}
}

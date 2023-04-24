package client

import (
	"log"

	"github.com/arvians-id/go-rabbitmq/todo/cmd/config"
	"github.com/arvians-id/go-rabbitmq/todo/pb"
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

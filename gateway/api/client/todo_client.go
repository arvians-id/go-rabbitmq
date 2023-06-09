package client

import (
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"google.golang.org/grpc"
	"log"
)

type TodoClient struct {
	Client pb.TodoServiceClient
}

func InitTodoClient(configuration config.Config) TodoClient {
	connection, err := grpc.Dial(configuration.Get("TODO_SERVICE_URL"), grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(config.NewGRPUnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalln(err)
	}

	return TodoClient{
		Client: pb.NewTodoServiceClient(connection),
	}
}

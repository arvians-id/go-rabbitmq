package client

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"google.golang.org/grpc"
	"log"
)

type TodoClient struct {
	Client pb.TodoServiceClient
}

func InitTodoClient(configuration config.Config) TodoClient {
	connection, err := grpc.Dial(configuration.Get("TODO_SERVICE_URL"), grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	return TodoClient{
		Client: pb.NewTodoServiceClient(connection),
	}
}

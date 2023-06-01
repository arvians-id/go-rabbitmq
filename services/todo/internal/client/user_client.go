package client

import (
	"log"

	"github.com/arvians-id/go-rabbitmq/todo/cmd/config"
	"github.com/arvians-id/go-rabbitmq/todo/pb"
	"google.golang.org/grpc"
)

type UserClient struct {
	Client pb.UserServiceClient
}

func InitUserClient(configuration config.Config) UserClient {
	connection, err := grpc.Dial(configuration.Get("USER_SERVICE_URL"), grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(config.NewGRPUnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalln(err)
	}

	return UserClient{
		Client: pb.NewUserServiceClient(connection),
	}
}

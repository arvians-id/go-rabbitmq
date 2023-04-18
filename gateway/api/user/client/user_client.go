package client

import (
	"log"

	"github.com/arvians-id/go-rabbitmq/gateway/api/user/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"google.golang.org/grpc"
)

type UserClient struct {
	Client pb.UserServiceClient
}

func InitUserClient(configuration config.Config) UserClient {
	connection, err := grpc.Dial(configuration.Get("USER_SERVICE_URL"), grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	return UserClient{
		Client: pb.NewUserServiceClient(connection),
	}
}

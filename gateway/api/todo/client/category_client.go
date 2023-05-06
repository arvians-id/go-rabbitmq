package client

import (
	"log"

	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"google.golang.org/grpc"
)

type CategoryClient struct {
	Client pb.CategoryServiceClient
}

func InitCategoryClient(configuration config.Config) CategoryClient {
	connection, err := grpc.Dial(configuration.Get("CATEGORY_SERVICE_URL"), grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(config.NewGRPUnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalln(err)
	}

	return CategoryClient{
		Client: pb.NewCategoryServiceClient(connection),
	}
}

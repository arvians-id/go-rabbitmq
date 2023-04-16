package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/arvians-id/go-rabbitmq/user/cmd/config"
	"github.com/arvians-id/go-rabbitmq/user/internal/repository"
	"github.com/arvians-id/go-rabbitmq/user/internal/usecase"
	"github.com/arvians-id/go-rabbitmq/user/pb"
	"google.golang.org/grpc"
)

func main() {
	configuration := config.New()
	db, err := config.NewPostgresSQL(configuration)
	if err != nil {
		log.Fatalln("Cannot connect to database", err)
	}

	userRepository := repository.NewUserRepository(db)
	userService := usecase.NewUserUsecase(userRepository)

	lis, err := net.Listen("tcp", configuration.Get("USER_SERVICE_URL"))
	if err != nil {
		log.Fatalln("Failed to listen", err)
	}

	port := strings.Split(configuration.Get("USER_SERVICE_URL"), ":")[1]
	fmt.Println("Category service is running on port", port)

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userService)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln("Failed to serving", err)
	}
}

package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/arvians-id/go-rabbitmq/todo/cmd/config"
	"github.com/arvians-id/go-rabbitmq/todo/internal/repository"
	"github.com/arvians-id/go-rabbitmq/todo/internal/usecase"
	"github.com/arvians-id/go-rabbitmq/todo/pb"
	"google.golang.org/grpc"
)

func main() {
	configuration := config.New(".env.dev")
	db, err := config.NewPostgresSQL(configuration)
	if err != nil {
		log.Fatalln("Cannot connect to database", err)
	}

	// Category Todo
	categoryTodoRepository := repository.NewCategoryTodoRepository(db)
	categoryTodoUsecase := usecase.NewCategoryTodoUsecase(categoryTodoRepository)

	// Todo
	todoRepository := repository.NewTodoRepository(db)
	todoService := usecase.NewTodoUsecase(todoRepository)

	lis, err := net.Listen("tcp", configuration.Get("TODO_SERVICE_URL"))
	if err != nil {
		log.Fatalln("Failed to listen", err)
	}

	port := strings.Split(configuration.Get("TODO_SERVICE_URL"), ":")[1]
	fmt.Println("Todo service is running on port", port)

	grpcServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(grpcServer, todoService)
	pb.RegisterCategoryTodoServiceServer(grpcServer, categoryTodoUsecase)
	
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln("Failed to serving", err)
	}
}

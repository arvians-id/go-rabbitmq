package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/arvians-id/go-rabbitmq/todo/cmd/config"
	"github.com/arvians-id/go-rabbitmq/todo/internal/client"
	"github.com/arvians-id/go-rabbitmq/todo/internal/repository"
	"github.com/arvians-id/go-rabbitmq/todo/internal/services"
	"github.com/arvians-id/go-rabbitmq/todo/internal/usecase"
	"github.com/arvians-id/go-rabbitmq/todo/pb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
)

func main() {
	tp, err := config.NewTracerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(ctx context.Context) {
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		err := tp.Shutdown(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	}(ctx)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	configuration := config.New(".env.dev")
	db, err := config.NewPostgresSQL(configuration)
	if err != nil {
		log.Fatalln("Cannot connect to database", err)
	}

	// User
	userClient := client.InitUserClient(configuration)
	userService := services.NewUserService(userClient)

	// Category Todo
	categoryTodoRepository := repository.NewCategoryTodoRepository(db)
	categoryTodoUsecase := usecase.NewCategoryTodoUsecase(categoryTodoRepository)

	// Todo
	todoRepository := repository.NewTodoRepository(db)
	todoService := usecase.NewTodoUsecase(userService, categoryTodoRepository, todoRepository)

	lis, err := net.Listen("tcp", configuration.Get("TODO_SERVICE_URL"))
	if err != nil {
		log.Fatalln("Failed to listen", err)
	}

	port := strings.Split(configuration.Get("TODO_SERVICE_URL"), ":")[1]
	fmt.Println("Todo service is running on port", port)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(config.NewGRPUnaryServerInterceptor()),
	)
	pb.RegisterTodoServiceServer(grpcServer, todoService)
	pb.RegisterCategoryTodoServiceServer(grpcServer, categoryTodoUsecase)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln("Failed to serving", err)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/arvians-id/go-rabbitmq/user/cmd/config"
	"github.com/arvians-id/go-rabbitmq/user/internal/repository"
	"github.com/arvians-id/go-rabbitmq/user/internal/usecase"
	"github.com/arvians-id/go-rabbitmq/user/pb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
)

func main() {
	// Init Config
	configuration := config.New()
	db, err := config.NewPostgresSQL(configuration)
	if err != nil {
		log.Fatalln("Cannot connect to database", err)
	}

	// Init Rabbit MQ
	conn, ch, err := config.InitRabbitMQ(configuration)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"mail",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Init Open Telementry Tracer
	tp, err := config.NewTracerProvider(configuration)
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

	// Init Server
	userRepository := repository.NewUserRepository(db)
	userService := usecase.NewUserUsecase(userRepository, ch)

	lis, err := net.Listen("tcp", configuration.Get("USER_SERVICE_URL"))
	if err != nil {
		log.Fatalln("Failed to listen", err)
	}

	port := strings.Split(configuration.Get("USER_SERVICE_URL"), ":")[1]
	fmt.Println("User service is running on port", port)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(config.NewGRPUnaryServerInterceptor()),
	)
	pb.RegisterUserServiceServer(grpcServer, userService)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln("Failed to serving", err)
	}
}

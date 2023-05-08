package main

import (
	"context"
	"fmt"
	"github.com/arvians-id/go-rabbitmq/category_todo/cmd/config"
	"github.com/arvians-id/go-rabbitmq/category_todo/internal/repository"
	"github.com/arvians-id/go-rabbitmq/category_todo/internal/usecase"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Init Config
	configuration := config.New()
	db, err := config.NewPostgresSQL(configuration)
	if err != nil {
		log.Fatalln("Cannot connect to database", err)
	}

	// Init RabbitMQ
	conn, ch, err := config.InitRabbitMQ(configuration)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	defer ch.Close()

	exchange := "category_todo_exchange"
	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("Cannot declare exchange", err)
	}

	fmt.Println("Category todo service is running")

	// Init Server
	// Category Todo Server
	ctx := context.Background()

	categoryTodoRepository := repository.NewCategoryTodoRepository(db)
	categoryTodoUsecase := usecase.NewCategoryTodoUsecase(categoryTodoRepository)

	go func() {
		err = categoryTodoUsecase.Delete(ctx, ch)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		err = categoryTodoUsecase.Create(ctx, ch)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
}

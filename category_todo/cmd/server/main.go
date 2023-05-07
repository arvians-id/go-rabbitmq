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
	"time"
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

	err = ch.ExchangeDeclare(
		"todos", // Nama exchange
		"topic", // Jenis exchange
		true,    // Durable
		false,   // Auto-deleted
		false,   // Internal
		false,   // No-wait
		nil,     // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

	q, err := ch.QueueDeclare(
		"",    // Nama antrian kosong (RabbitMQ akan memberikan nama antrian yang unik)
		false, // Non-durable
		false, // Non-autodelete
		true,  // Exclusive (antrian hanya dapat digunakan oleh satu subscriber)
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Init Server
	// Category Todo Server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	categoryTodoRepository := repository.NewCategoryTodoRepository(db)
	categoryTodoUsecase := usecase.NewCategoryTodoUsecase(categoryTodoRepository)
	_, err = categoryTodoUsecase.Create(ctx, ch, q)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Category todo service is running")

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
}

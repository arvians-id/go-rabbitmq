package main

import (
	"context"
	"fmt"
	"github.com/arvians-id/go-rabbitmq/category_todo/cmd/config"
	"github.com/arvians-id/go-rabbitmq/category_todo/internal"
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
		log.Fatalln("Cannot connect to RabbitMQ", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatalln("Cannot set Qos", err)
	}

	fmt.Println("Category todo service is running")

	// Init Server
	// Category Todo Server
	ctx := context.Background()
	internal.NewApp(ctx, ch, db)

	defer func() {
		conn.Close()
		ch.Close()
	}()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
}

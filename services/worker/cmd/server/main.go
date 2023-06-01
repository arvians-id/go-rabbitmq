package main

import (
	"fmt"
	"github.com/arvians-id/go-rabbitmq/worker/cmd/config"
	"github.com/arvians-id/go-rabbitmq/worker/internal"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Init Config
	configuration := config.New()

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

	fmt.Println("Worker message service is running")

	// Init Server
	internal.NewApp(configuration, ch)

	defer func() {
		conn.Close()
		ch.Close()
	}()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
}

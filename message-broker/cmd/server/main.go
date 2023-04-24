package main

import (
	"bytes"
	"github.com/arvians-id/go-rabbitmq/message-broker/cmd/config"
	"log"
	"time"
)

func main() {
	conn, ch, err := config.InitRabbitMQ()
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	message, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	var forever chan struct{}

	go func() {
		for d := range message {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

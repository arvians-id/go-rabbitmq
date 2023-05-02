package main

import (
	"encoding/json"
	"github.com/arvians-id/go-rabbitmq/message-broker/cmd/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configuration := config.New(".env.dev")
	conn, ch, err := config.InitRabbitMQ(configuration)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare(
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

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

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

	go func() {
		type Message struct {
			ToEmail string
			Message string
		}
		for d := range message {
			var mailMessage Message
			err := json.Unmarshal(d.Body, &mailMessage)
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("Received a message: %s", mailMessage.ToEmail)
			err = config.SendMail(configuration, mailMessage.ToEmail, "Test Subject Mail", mailMessage.Message)
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("Mail Sended")
			d.Ack(false)
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
}

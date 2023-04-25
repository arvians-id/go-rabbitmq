package main

import (
	"log"
	"time"

	"github.com/arvians-id/go-rabbitmq/message-broker/cmd/config"
)

func main() {
	configuration := config.New(".env.dev")
	conn, ch, err := config.InitRabbitMQ()
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
			err := config.SendMail(configuration, string(d.Body), "Test Subject Mail", "Hallo bang hehe")
			if err != nil {
				log.Fatalln(err)
			}
			time.Sleep(5 * time.Second)
			log.Printf("mail Sended")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

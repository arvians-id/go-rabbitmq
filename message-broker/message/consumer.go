package message

import (
	"github.com/rabbitmq/amqp091-go"
	"log"
)

func Consume(ch *amqp091.Channel, q amqp091.Queue) {
	message, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}

	var forever chan struct{}

	go func() {
		for d := range message {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

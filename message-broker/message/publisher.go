package message

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func Publish(ch *amqp091.Channel, q amqp091.Queue) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	body := "Hello world!"
	err := ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return err
	}
	log.Printf(" [x] Sent %s\n", body)

	return nil
}

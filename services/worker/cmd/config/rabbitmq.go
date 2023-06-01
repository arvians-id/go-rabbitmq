package config

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ(configuration Config) (*amqp091.Connection, *amqp091.Channel, error) {
	mqUser := configuration.Get("MQ_USER")
	mqPassword := configuration.Get("MQ_PASSWORD")
	mqHost := configuration.Get("MQ_HOST")
	mqPort := configuration.Get("MQ_PORT")

	// connect to rabbitmq
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", mqUser, mqPassword, mqHost, mqPort)
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

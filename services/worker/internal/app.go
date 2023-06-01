package internal

import (
	"github.com/arvians-id/go-rabbitmq/worker/cmd/config"
	"github.com/arvians-id/go-rabbitmq/worker/internal/usecase"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

func NewApp(configuration config.Config, channel *amqp091.Channel) {
	mailUsecase := usecase.NewMailUsecase(configuration)

	go func() {
		err := mailUsecase.SendMail(channel)
		if err != nil {
			log.Fatalln(err)
		}
	}()
}

package usecase

import (
	"github.com/arvians-id/go-rabbitmq/worker/cmd/config"
	"github.com/arvians-id/go-rabbitmq/worker/internal/model"
	"github.com/goccy/go-json"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

type MailUsecase struct {
	Configuration config.Config
}

func NewMailUsecase(configuration config.Config) *MailUsecase {
	return &MailUsecase{
		Configuration: configuration,
	}
}

func (usecase *MailUsecase) SendMail(channel *amqp091.Channel) error {
	q, err := channel.QueueDeclare(
		"mail",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	message, err := channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var errs error

	go func() {
		defer wg.Done()
		for d := range message {
			var mailMessage model.Message
			err := json.Unmarshal(d.Body, &mailMessage)
			if err != nil {
				errs = err
				return
			}

			log.Printf("Received a message: %s", mailMessage.ToEmail)
			err = config.SendMail(usecase.Configuration, mailMessage.ToEmail, "Test Subject Mail", mailMessage.Message)
			if err != nil {
				errs = err
				return
			}

			log.Printf("Mail Sended")
			d.Ack(false)
		}
	}()
	wg.Wait()
	if errs != nil {
		return errs
	}

	return nil
}

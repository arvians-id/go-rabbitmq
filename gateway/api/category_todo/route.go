package category_todo

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/category_todo/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
)

func NewCategoryTodoRoute(c fiber.Router, channel *amqp091.Channel) error {
	exchangeName := "category_todo_exchange"
	err := channel.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	categoryTodoHandler := handler.NewCategoryTodoHandler(channel)

	c.Delete("/category-todo", categoryTodoHandler.Delete)

	return nil
}

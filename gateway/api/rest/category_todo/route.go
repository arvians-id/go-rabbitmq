package category_todo

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/middleware"
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/category_todo/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

func NewCategoryTodoRoute(c fiber.Router, channel *amqp091.Channel) {
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
		log.Fatalln(err)
	}

	categoryTodoHandler := handler.NewCategoryTodoHandler(channel)

	c.Delete("/category-todo", middleware.NewJWTMiddleware(), categoryTodoHandler.Delete)
}

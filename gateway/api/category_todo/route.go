package category_todo

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/category_todo/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
)

func NewCategoryTodoRoute(c fiber.Router, channel *amqp091.Channel) {
	categoryTodoHandler := handler.NewCategoryTodoHandler(channel)

	c.Delete("/category-todo", categoryTodoHandler.Delete)
}

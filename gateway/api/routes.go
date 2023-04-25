package api

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo"
	"github.com/arvians-id/go-rabbitmq/gateway/api/user"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
)

func NewRoutes(c *fiber.App, channel *amqp091.Channel, configuration config.Config) {
	apiGroup := c.Group("/api")
	user.NewUserRoute(apiGroup, configuration)
	todo.NewTodoRoute(apiGroup, channel, configuration)
	todo.NewCategoryTodoRoute(apiGroup, configuration)
}

package auth

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/auth/handler"
	"github.com/arvians-id/go-rabbitmq/gateway/api/services"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
)

func NewAuthRoute(c fiber.Router, configuration config.Config, channel *amqp091.Channel) {
	userClient := client.InitUserClient(configuration)
	userService := services.NewUserService(userClient, channel)
	userHandler := handler.NewLoginHandler(userService)

	c.Post("/register", userHandler.Register)
	c.Post("/login", userHandler.Login)
}

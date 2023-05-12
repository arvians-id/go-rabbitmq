package auth

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/auth/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/auth/handler"
	"github.com/arvians-id/go-rabbitmq/gateway/api/auth/services"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
)

func NewAuthRoute(c fiber.Router, configuration config.Config, channel *amqp091.Channel) {
	userClient := client.InitUserClient(configuration)
	userService := services.NewUserService(userClient)
	userHandler := handler.NewLoginHandler(userService, channel)

	c.Post("/register", userHandler.Register)
	c.Post("/login", userHandler.Login)
}

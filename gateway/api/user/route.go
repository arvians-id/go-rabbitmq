package user

import (
	"errors"
	"github.com/arvians-id/go-rabbitmq/gateway/api/user/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/user/handler"
	"github.com/arvians-id/go-rabbitmq/gateway/api/user/services"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

func NewUserRoute(c fiber.Router, configuration config.Config, channel *amqp091.Channel) {
	userClient := client.InitUserClient(configuration)
	userService := services.NewUserService(userClient)
	userHandler := handler.NewUserHandler(userService, channel)

	c.Get("/users", userHandler.FindAll)
	c.Get("/users/:id", userHandler.FindByID)
	c.Post("/users", userHandler.Create)
	c.Patch("/users/:id", limiter.New(limiter.Config{
		Max:        5,
		Expiration: 3 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return response.ReturnError(c, fiber.StatusTooManyRequests, errors.New("too many requests, please try again later after 3 minute"))
		},
	}), userHandler.Update)
	c.Delete("/users/:id", userHandler.Delete)
}

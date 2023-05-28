package user

import (
	"errors"
	"github.com/arvians-id/go-rabbitmq/gateway/api/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/middleware"
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/user/handler"
	"github.com/arvians-id/go-rabbitmq/gateway/api/services"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

func NewUserRoute(c fiber.Router, configuration config.Config, channel *amqp091.Channel) services.UserServiceContract {
	userClient := client.InitUserClient(configuration)
	userService := services.NewUserService(userClient, channel)
	userHandler := handler.NewUserHandler(userService)

	routeUser := c.Group("/users", middleware.NewJWTMiddleware())
	routeUser.Get("/", userHandler.FindAll)
	routeUser.Get("/:id", userHandler.FindByID)
	routeUser.Post("/", userHandler.Create)
	routeUser.Patch("/:id", limiter.New(limiter.Config{
		Max:        5,
		Expiration: 3 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return response.ReturnError(c, fiber.StatusTooManyRequests, errors.New("too many requests, please try again later after 3 minute"))
		},
	}), userHandler.Update)
	routeUser.Delete("/:id", userHandler.Delete)

	return userService
}

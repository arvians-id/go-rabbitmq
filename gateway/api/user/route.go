package user

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/user/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/user/handler"
	"github.com/arvians-id/go-rabbitmq/gateway/api/user/services"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/gofiber/fiber/v2"
)

func NewUserRoute(c fiber.Router, configuration config.Config) {
	userClient := client.InitUserClient(configuration)
	userService := services.NewUserService(userClient)
	userHandler := handler.NewUserHandler(userService)

	c.Get("/users", userHandler.FindAll)
	c.Get("/users/:id", userHandler.FindByID)
	c.Post("/users", userHandler.Create)
	c.Patch("/users/:id", userHandler.Update)
	c.Delete("/users/:id", userHandler.Delete)
}

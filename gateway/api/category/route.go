package category

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/category/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/category/handler"
	"github.com/arvians-id/go-rabbitmq/gateway/api/category/services"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/gofiber/fiber/v2"
)

func NewCategoryRoute(c fiber.Router, configuration config.Config) {
	categoryClient := client.InitCategoryClient(configuration)
	categoryService := services.NewCategoryService(categoryClient)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	c.Get("/categories", categoryHandler.FindAll)
	c.Get("/categories/:id", categoryHandler.FindByID)
	c.Post("/categories", categoryHandler.Create)
	c.Delete("/categories/:id", categoryHandler.Delete)
}

package category

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/middleware"
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/category/handler"
	"github.com/arvians-id/go-rabbitmq/gateway/api/services"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/gofiber/fiber/v2"
)

func NewCategoryRoute(c fiber.Router, configuration config.Config) services.CategoryServiceContract {
	categoryClient := client.InitCategoryClient(configuration)
	categoryService := services.NewCategoryService(categoryClient)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	routeCategory := c.Group("/categories", middleware.NewJWTMiddleware())
	routeCategory.Get("/", categoryHandler.FindAll)
	routeCategory.Get("/:id", categoryHandler.FindByID)
	routeCategory.Post("/", categoryHandler.Create)
	routeCategory.Delete("/:id", categoryHandler.Delete)
	routeCategory.Get("/:todoId/todo", categoryHandler.FindAllByTodoID)

	return categoryService
}

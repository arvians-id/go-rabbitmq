package todo

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/handler"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/services"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/gofiber/fiber/v2"
)

func NewTodoRoute(c fiber.Router, configuration config.Config) {
	todoClient := client.InitTodoClient(configuration)
	todoService := services.NewTodoService(todoClient)
	todoHandler := handler.NewTodoHandler(todoService)

	c.Get("/todos", todoHandler.FindAll)
	c.Get("/todos/:id", todoHandler.FindByID)
	c.Post("/todos", todoHandler.Create)
	c.Patch("/todos/:id", todoHandler.Update)
	c.Delete("/todos/:id", todoHandler.Delete)
}

func NewCategoryTodoRoute(c fiber.Router, configuration config.Config) {
	categoryTodoClient := client.InitCategoryTodoClient(configuration)
	categoryTodoService := services.NewCategoryTodoService(categoryTodoClient)
	categoryTodoHandler := handler.NewCategoryTodoHandler(categoryTodoService)

	c.Get("/category-todos", categoryTodoHandler.FindAll)
	c.Get("/category-todos/:id", categoryTodoHandler.FindByID)
	c.Post("/category-todos", categoryTodoHandler.Create)
	c.Delete("/category-todos/:id", categoryTodoHandler.Delete)
}

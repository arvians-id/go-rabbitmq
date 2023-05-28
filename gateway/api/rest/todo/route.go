package todo

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/client"
	"github.com/arvians-id/go-rabbitmq/gateway/api/middleware"
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/todo/handler"
	"github.com/arvians-id/go-rabbitmq/gateway/api/services"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func NewTodoRoute(c fiber.Router, configuration config.Config, redisClient *redis.Client) services.TodoServiceContract {
	categoryClient := client.InitCategoryClient(configuration)
	categoryService := services.NewCategoryService(categoryClient)

	todoClient := client.InitTodoClient(configuration)
	todoService := services.NewTodoServiceCache(todoClient, redisClient, categoryService)
	todoHandler := handler.NewTodoHandler(todoService, categoryService)

	c.Get("/display-todos", middleware.NewJWTMiddleware(), todoHandler.DisplayTodoCategoryList)
	routeTodo := c.Group("/todos", middleware.NewJWTMiddleware())
	routeTodo.Get("/", todoHandler.FindAll)
	routeTodo.Get("/:id", todoHandler.FindByID)
	routeTodo.Post("/", todoHandler.Create)
	routeTodo.Patch("/:id", todoHandler.Update)
	routeTodo.Delete("/:id", todoHandler.Delete)

	return todoService
}

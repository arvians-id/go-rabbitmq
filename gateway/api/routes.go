package api

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo"
	"github.com/arvians-id/go-rabbitmq/gateway/api/user"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func NewRoutes(c *fiber.App, configuration config.Config, redisClient *redis.Client) {
	apiGroup := c.Group("/api")
	user.NewUserRoute(apiGroup, configuration)
	todo.NewTodoRoute(apiGroup, configuration, redisClient)
	todo.NewCategoryTodoRoute(apiGroup, configuration)
}

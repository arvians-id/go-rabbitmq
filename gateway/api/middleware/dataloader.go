package middleware

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/gql/model"
	"github.com/arvians-id/go-rabbitmq/gateway/api/services"
	"github.com/gofiber/fiber/v2"
)

type Loaders struct {
	UserServiceFindByIDs         model.UserLoader
	TodoServiceFindByUserIDs     model.TodoLoader
	CategoryServiceFindByTodoIDs model.CategoryLoader
}

func DataLoaders(
	userService services.UserServiceContract,
	todoService services.TodoServiceContract,
	categoryService services.CategoryServiceContract,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		loaders := Loaders{
			UserServiceFindByIDs:         model.UserServiceFindByIDs(c.Context(), userService),
			TodoServiceFindByUserIDs:     model.TodoServiceFindByUserIDs(c.Context(), todoService),
			CategoryServiceFindByTodoIDs: model.CategoryServiceFindByTodoIDs(c.Context(), categoryService),
		}

		c.Locals("loaders", &loaders)
		return c.Next()
	}
}

func GetLoaders(ctx context.Context) *Loaders {
	return ctx.Value("loaders").(*Loaders)
}

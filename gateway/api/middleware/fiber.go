package middleware

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

func ExposeFiberContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("fiberContext", c)
		return c.Next()
	}
}

func FiberContext(ctx context.Context) *fiber.Ctx {
	return ctx.Value("fiberContext").(*fiber.Ctx)
}

package middleware

import (
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func NewJWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return response.ReturnError(ctx, fiber.StatusUnauthorized, err)
		},
		SuccessHandler: func(ctx *fiber.Ctx) error {
			userContext := ctx.Locals("user").(*jwt.Token)
			userClaims := userContext.Claims.(jwt.MapClaims)

			ctx.Locals("id", userClaims["id"])
			ctx.Locals("email", userClaims["email"])
			return ctx.Next()
		},
	})
}

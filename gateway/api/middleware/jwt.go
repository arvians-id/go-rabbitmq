package middleware

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"strings"
)

var SecretKey = []byte("secret")

func NewJWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: SecretKey,
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

func NewJWTMiddlewareGraphQL(ctx context.Context, obj interface{}, next graphql.Resolver, isLoggedIn bool) (interface{}, error) {
	fiberCtx := FiberContext(ctx)

	if !isLoggedIn {
		return next(ctx)
	}

	authorizationHeader := fiberCtx.Get("Authorization")
	if !strings.Contains(authorizationHeader, "Bearer") {
		return nil, gqlerror.Errorf("invalid authorization header")
	}

	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, gqlerror.Errorf("Invalid signing method")
		}
		return SecretKey, nil
	})

	if err != nil {
		return nil, gqlerror.Errorf("Invalid JWT token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx = context.WithValue(ctx, "id", claims["id"])
		ctx = context.WithValue(ctx, "email", claims["email"])
	} else {
		return nil, gqlerror.Errorf("Invalid JWT token")
	}

	return next(ctx)
}

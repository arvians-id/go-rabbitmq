package response

import (
	"github.com/gofiber/fiber/v2"
)

var (
	GrpcErrorNotFound = "rpc error: code = Unknown desc = sql: no rows in result set"
	//GrpcErrorNotFound = "rpc error: code = Unknown desc = sql: no rows in result set"
)

func ReturnError(c *fiber.Ctx, code int, err error) error {
	return c.Status(code).JSON(ApiResponse{
		Code:   code,
		Status: err.Error(),
		Data:   nil,
	})
}

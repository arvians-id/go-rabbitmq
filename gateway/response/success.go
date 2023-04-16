package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ReturnSuccessOK(c *fiber.Ctx, status string, data interface{}) error {
	c.Status(http.StatusOK)
	return c.JSON(WebResponse{
		Code:   http.StatusOK,
		Status: status,
		Data:   data,
	})
}

func ReturnSuccessCreated(c *fiber.Ctx, status string, data interface{}) error {
	c.Status(http.StatusCreated)
	return c.JSON(WebResponse{
		Code:   http.StatusCreated,
		Status: status,
		Data:   data,
	})
}

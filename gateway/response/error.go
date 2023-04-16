package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ReturnErrorNotFound(c *fiber.Ctx, err error) error {
	c.Status(http.StatusNotFound)
	return c.JSON(WebResponse{
		Code:   http.StatusNotFound,
		Status: "data not found",
		Data:   nil,
	})
}

func ReturnErrorInternalServerError(c *fiber.Ctx, err error) error {
	c.Status(http.StatusInternalServerError)
	return c.JSON(WebResponse{
		Code:   http.StatusInternalServerError,
		Status: err.Error(),
		Data:   nil,
	})
}

func ReturnErrorBadRequest(c *fiber.Ctx, err error) error {
	return c.Status(http.StatusBadRequest).JSON(WebResponse{
		Code:   http.StatusBadRequest,
		Status: err.Error(),
		Data:   nil,
	})
}

func ReturnErrorUnauthorized(c *fiber.Ctx, err error) error {
	c.Status(http.StatusUnauthorized)
	return c.JSON(WebResponse{
		Code:   http.StatusUnauthorized,
		Status: err.Error(),
		Data:   nil,
	})
}

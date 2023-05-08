package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/dto"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
)

type CategoryTodoHandler struct {
	RabbitMQ *amqp091.Channel
}

func NewCategoryTodoHandler(rabbitMQ *amqp091.Channel) CategoryTodoHandler {
	return CategoryTodoHandler{
		RabbitMQ: rabbitMQ,
	}
}

func (handler *CategoryTodoHandler) Delete(c *fiber.Ctx) error {
	var categoryTodo dto.CategoryTodo
	err := c.BodyParser(&categoryTodo)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	exchangeName := "category_todo_exchange"
	err = handler.RabbitMQ.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	data, err := json.Marshal(categoryTodo)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = handler.RabbitMQ.PublishWithContext(
		c.Context(),
		exchangeName,
		"category_todo.deleted",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "deleted", nil)
}

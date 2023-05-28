package handler

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/category_todo/dto"
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

	err = handler.publish(c.Context(), "category_todo.deleted", categoryTodo)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "deleted", nil)
}

func (handler *CategoryTodoHandler) publish(ctx context.Context, key string, data interface{}) error {
	marshaled, err := json.Marshal(data)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = handler.RabbitMQ.PublishWithContext(
		ctx,
		"category_todo_exchange",
		key,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        marshaled,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/request"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/services"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TodoHandler struct {
	TodoService services.TodoService
	RabbitMQ    *amqp091.Channel
}

func NewTodoHandler(todoService services.TodoService, rabbitMQ *amqp091.Channel) TodoHandler {
	return TodoHandler{
		TodoService: todoService,
		RabbitMQ:    rabbitMQ,
	}
}

func (handler *TodoHandler) FindAll(c *fiber.Ctx) error {
	users, err := handler.TodoService.FindAll(c.Context(), new(emptypb.Empty))
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessOK(c, "OK", users)
}

func (handler *TodoHandler) FindByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	todo, err := handler.TodoService.FindByID(c.Context(), &pb.GetTodoByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessOK(c, "OK", todo)
}

func (handler *TodoHandler) Create(c *fiber.Ctx) error {
	var todoRequest request.TodoCreateRequest
	err := c.BodyParser(&todoRequest)
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	err = helper.ValidateStruct(todoRequest)
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	todoCreated, err := handler.TodoService.Create(c.Context(), &pb.CreateTodoRequest{
		Title:          todoRequest.Title,
		Description:    todoRequest.Description,
		UserId:         todoRequest.UserId,
		CategoryTodoId: todoRequest.CategoryTodoId,
	})
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	// Send Email To Queue
	err = handler.RabbitMQ.PublishWithContext(
		c.Context(),
		"",
		"mail",
		false,
		false,
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "text/plain",
			Body:         []byte("widdyarfiansyah00@gmail.com"),
		},
	)
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessCreated(c, "created", todoCreated)
}

func (handler *TodoHandler) Update(c *fiber.Ctx) error {
	var todoRequest request.TodoUpdateRequest
	err := c.BodyParser(&todoRequest)
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	err = helper.ValidateStruct(todoRequest)
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	todoUpdated, err := handler.TodoService.Update(c.Context(), &pb.UpdateTodoRequest{
		Id:          int64(id),
		Title:       todoRequest.Title,
		Description: todoRequest.Description,
		IsDone:      todoRequest.IsDone,
		UserId:      todoRequest.UserId,
	})
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessOK(c, "updated", todoUpdated)
}

func (handler *TodoHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	_, err = handler.TodoService.Delete(c.Context(), &pb.GetTodoByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessOK(c, "deleted", nil)
}

package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/request"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/services"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryTodoHandler struct {
	CategoryTodoService services.CategoryTodoServiceContract
}

func NewCategoryTodoHandler(categoryTodoService services.CategoryTodoServiceContract) CategoryTodoHandler {
	return CategoryTodoHandler{
		CategoryTodoService: categoryTodoService,
	}
}

func (handler *CategoryTodoHandler) FindAll(c *fiber.Ctx) error {
	categoryTodos, err := handler.CategoryTodoService.FindAll(c.Context(), new(emptypb.Empty))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "OK", categoryTodos)
}

func (handler *CategoryTodoHandler) FindByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	categoryTodo, err := handler.CategoryTodoService.FindByID(c.Context(), &pb.GetCategoryTodoByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "OK", categoryTodo)
}

func (handler *CategoryTodoHandler) Create(c *fiber.Ctx) error {
	var categoryTodoRequest request.CategoryTodoCreateRequest
	err := c.BodyParser(&categoryTodoRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = helper.ValidateStruct(categoryTodoRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	categoryTodoCreated, err := handler.CategoryTodoService.Create(c.Context(), &pb.CreateCategoryTodoRequest{
		Name: categoryTodoRequest.Name,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusCreated, "created", categoryTodoCreated)
}

func (handler *CategoryTodoHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err = handler.CategoryTodoService.Delete(c.Context(), &pb.GetCategoryTodoByIDRequest{
		Id: int64(id),
	})

	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "deleted", nil)
}

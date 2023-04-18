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
	CategoryTodoService services.CategoryTodoService
}

func NewCategoryTodoHandler(categoryTodoService services.CategoryTodoService) CategoryTodoHandler {
	return CategoryTodoHandler{
		CategoryTodoService: categoryTodoService,
	}
}

func (handler *CategoryTodoHandler) FindAll(c *fiber.Ctx) error {
	categoryTodos, err := handler.CategoryTodoService.FindAll(c.Context(), new(emptypb.Empty))
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessOK(c, "OK", categoryTodos)
}

func (handler *CategoryTodoHandler) FindByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	categoryTodo, err := handler.CategoryTodoService.FindByID(c.Context(), &pb.GetCategoryTodoByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessOK(c, "OK", categoryTodo)
}

func (handler *CategoryTodoHandler) Create(c *fiber.Ctx) error {
	var categoryTodoRequest request.CategoryTodoCreateRequest
	err := c.BodyParser(&categoryTodoRequest)
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	err = helper.ValidateStruct(categoryTodoRequest)
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	categoryTodoCreated, err := handler.CategoryTodoService.Create(c.Context(), &pb.CreateCategoryTodoRequest{
		Name: categoryTodoRequest.Name,
	})
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessCreated(c, "created", categoryTodoCreated)
}

func (handler *CategoryTodoHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	_, err = handler.CategoryTodoService.Delete(c.Context(), &pb.GetCategoryTodoByIDRequest{
		Id: int64(id),
	})

	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessOK(c, "deleted", nil)
}

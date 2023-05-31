package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/todo/request"
	"github.com/arvians-id/go-rabbitmq/gateway/api/services"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
)

type TodoHandler struct {
	TodoService     services.TodoServiceContract
	CategoryService services.CategoryServiceContract
}

func NewTodoHandler(todoService services.TodoServiceContract, categoryService services.CategoryServiceContract) TodoHandler {
	return TodoHandler{
		TodoService:     todoService,
		CategoryService: categoryService,
	}
}

func (handler *TodoHandler) DisplayTodoCategoryList(c *fiber.Ctx) error {
	data, code, err := handler.TodoService.DisplayTodoCategoryList(c.Context())
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "OK", data)
}

func (handler *TodoHandler) FindAll(c *fiber.Ctx) error {
	todos, code, err := handler.TodoService.FindAll(c.Context())
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "OK", todos.GetTodos())
}

func (handler *TodoHandler) FindByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	todo, code, err := handler.TodoService.FindByID(c.Context(), &pb.GetTodoByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "OK", todo.GetTodo())
}

func (handler *TodoHandler) Create(c *fiber.Ctx) error {
	var todoRequest request.TodoCreateRequest
	err := c.BodyParser(&todoRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = helper.ValidateStruct(todoRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	todoCreated, code, err := handler.TodoService.Create(c.Context(), &pb.CreateTodoRequest{
		Title:       todoRequest.Title,
		Description: todoRequest.Description,
		UserId:      todoRequest.UserId,
		CategoryId:  todoRequest.Categories,
	})
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "created", todoCreated.GetTodo())
}

func (handler *TodoHandler) Update(c *fiber.Ctx) error {
	var todoRequest request.TodoUpdateRequest
	err := c.BodyParser(&todoRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = helper.ValidateStruct(todoRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	todoUpdated, code, err := handler.TodoService.Update(c.Context(), &pb.UpdateTodoRequest{
		Id:          int64(id),
		Title:       todoRequest.Title,
		Description: todoRequest.Description,
		IsDone:      &todoRequest.IsDone,
		UserId:      todoRequest.UserId,
		CategoryId:  todoRequest.Categories,
	})
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "updated", todoUpdated.GetTodo())
}

func (handler *TodoHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	code, err := handler.TodoService.Delete(c.Context(), &pb.GetTodoByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "deleted", nil)
}

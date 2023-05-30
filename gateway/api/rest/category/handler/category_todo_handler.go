package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/category/request"
	"github.com/arvians-id/go-rabbitmq/gateway/api/services"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	CategoryService services.CategoryServiceContract
}

func NewCategoryHandler(categoryService services.CategoryServiceContract) CategoryHandler {
	return CategoryHandler{
		CategoryService: categoryService,
	}
}

func (handler *CategoryHandler) FindAll(c *fiber.Ctx) error {
	categories, code, err := handler.CategoryService.FindAll(c.Context())
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "OK", categories.GetCategories())
}

func (handler *CategoryHandler) FindAllByTodoID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("todoId")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	categories, code, err := handler.CategoryService.FindAllByTodoID(c.Context(), &pb.GetCategoryByTodoIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "OK", categories.GetCategories())
}

func (handler *CategoryHandler) FindByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	category, code, err := handler.CategoryService.FindByID(c.Context(), &pb.GetCategoryByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "OK", category.GetCategory())
}

func (handler *CategoryHandler) Create(c *fiber.Ctx) error {
	var categoryRequest request.CategoryCreateRequest
	err := c.BodyParser(&categoryRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = helper.ValidateStruct(categoryRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	categoryCreated, code, err := handler.CategoryService.Create(c.Context(), &pb.CreateCategoryRequest{
		Name: categoryRequest.Name,
	})
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "created", categoryCreated.GetCategory())
}

func (handler *CategoryHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	code, err := handler.CategoryService.Delete(c.Context(), &pb.GetCategoryByIDRequest{
		Id: int64(id),
	})

	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "deleted", nil)
}

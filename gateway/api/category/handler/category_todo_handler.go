package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/category/request"
	"github.com/arvians-id/go-rabbitmq/gateway/api/category/services"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/emptypb"
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
	categorys, err := handler.CategoryService.FindAll(c.Context(), new(emptypb.Empty))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "OK", categorys.GetCategories())
}

func (handler *CategoryHandler) FindByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	category, err := handler.CategoryService.FindByID(c.Context(), &pb.GetCategoryByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "OK", category.GetCategory())
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

	categoryCreated, err := handler.CategoryService.Create(c.Context(), &pb.CreateCategoryRequest{
		Name: categoryRequest.Name,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusCreated, "created", categoryCreated.GetCategory())
}

func (handler *CategoryHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err = handler.CategoryService.Delete(c.Context(), &pb.GetCategoryByIDRequest{
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

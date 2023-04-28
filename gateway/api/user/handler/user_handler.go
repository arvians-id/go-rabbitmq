package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/user/request"
	"github.com/arvians-id/go-rabbitmq/gateway/api/user/services"

	"github.com/arvians-id/go-rabbitmq/gateway/api/user/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

const name = "gateway"

type UserHandler struct {
	UserService services.UserServiceContract
}

func NewUserHandler(userService services.UserServiceContract) UserHandler {
	return UserHandler{
		UserService: userService,
	}
}

func (handler *UserHandler) FindAll(c *fiber.Ctx) error {
	users, err := handler.UserService.FindAll(c.Context(), new(emptypb.Empty))
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessOK(c, "OK", users)
}

func (handler *UserHandler) FindByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	user, err := handler.UserService.FindByID(c.Context(), &pb.GetUserByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessOK(c, "OK", user)
}

func (handler *UserHandler) Create(c *fiber.Ctx) error {
	var userRequest request.UserCreateRequest
	err := c.BodyParser(&userRequest)
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	err = helper.ValidateStruct(userRequest)
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	userCreated, err := handler.UserService.Create(c.Context(), &pb.CreateUserRequest{
		Name:  userRequest.Name,
		Email: userRequest.Email,
	})
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessCreated(c, "created", userCreated)
}

func (handler *UserHandler) Update(c *fiber.Ctx) error {
	var userRequest request.UserUpdateRequest
	err := c.BodyParser(&userRequest)
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	errValidate := helper.ValidateStruct(userRequest)
	if errValidate != nil {
		return response.ReturnErrorBadRequest(c, errValidate)
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	userUpdated, err := handler.UserService.Update(c.Context(), &pb.UpdateUserRequest{
		Id:   int64(id),
		Name: userRequest.Name,
	})

	return response.ReturnSuccessOK(c, "updated", userUpdated)
}

func (handler *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	_, err = handler.UserService.Delete(c.Context(), &pb.GetUserByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return err
	}

	return response.ReturnSuccessOK(c, "deleted", nil)
}

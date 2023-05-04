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
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "OK", users.GetUsers())
}

func (handler *UserHandler) FindByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := handler.UserService.FindByID(c.Context(), &pb.GetUserByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "OK", user.GetUser())
}

func (handler *UserHandler) Create(c *fiber.Ctx) error {
	var userRequest request.UserCreateRequest
	err := c.BodyParser(&userRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = helper.ValidateStruct(userRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userCreated, err := handler.UserService.Create(c.Context(), &pb.CreateUserRequest{
		Name:  userRequest.Name,
		Email: userRequest.Email,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusCreated, "created", userCreated.GetUser())
}

func (handler *UserHandler) Update(c *fiber.Ctx) error {
	var userRequest request.UserUpdateRequest
	err := c.BodyParser(&userRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = helper.ValidateStruct(userRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userUpdated, err := handler.UserService.Update(c.Context(), &pb.UpdateUserRequest{
		Id:   int64(id),
		Name: userRequest.Name,
	})
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "updated", userUpdated.GetUser())
}

func (handler *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err = handler.UserService.Delete(c.Context(), &pb.GetUserByIDRequest{
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

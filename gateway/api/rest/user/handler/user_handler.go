package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/user/request"
	"github.com/arvians-id/go-rabbitmq/gateway/api/services"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserService services.UserServiceContract
}

func NewUserHandler(userService services.UserServiceContract) UserHandler {
	return UserHandler{
		UserService: userService,
	}
}

func (handler *UserHandler) FindAll(c *fiber.Ctx) error {
	userIds := c.Query("ids")
	if userIds != "" {
		ids, err := helper.ConvertStringToBulkInt64(userIds)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		users, code, err := handler.UserService.FindByIDs(c.Context(), &pb.GetUserByIDsRequest{
			Ids: ids,
		})
		if err != nil {
			return fiber.NewError(code, err.Error())
		}

		return response.ReturnSuccess(c, code, "OK", users.GetUsers())
	}
	users, code, err := handler.UserService.FindAll(c.Context())
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "OK", users.GetUsers())
}

func (handler *UserHandler) FindByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, code, err := handler.UserService.FindByID(c.Context(), &pb.GetUserByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "OK", user.GetUser())
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

	userCreated, code, err := handler.UserService.Create(c.Context(), &pb.CreateUserRequest{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: userRequest.Password,
	})
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "created", &pb.User{
		Id:        userCreated.GetUser().GetId(),
		Name:      userCreated.GetUser().GetName(),
		Email:     userCreated.GetUser().GetEmail(),
		CreatedAt: userCreated.GetUser().GetCreatedAt(),
		UpdatedAt: userCreated.GetUser().GetUpdatedAt(),
	})
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

	userUpdated, code, err := handler.UserService.Update(c.Context(), &pb.UpdateUserRequest{
		Id:       int64(id),
		Name:     userRequest.Name,
		Password: userRequest.Password,
	})
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "updated", &pb.User{
		Id:        userUpdated.GetUser().GetId(),
		Name:      userUpdated.GetUser().GetName(),
		Email:     userUpdated.GetUser().GetEmail(),
		CreatedAt: userUpdated.GetUser().GetCreatedAt(),
		UpdatedAt: userUpdated.GetUser().GetUpdatedAt(),
	})
}

func (handler *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	code, err := handler.UserService.Delete(c.Context(), &pb.GetUserByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "deleted", nil)
}

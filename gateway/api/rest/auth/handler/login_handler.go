package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/rest/auth/request"
	"github.com/arvians-id/go-rabbitmq/gateway/api/services"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
)

type LoginHandler struct {
	UserService services.UserServiceContract
}

func NewLoginHandler(userService services.UserServiceContract) *LoginHandler {
	return &LoginHandler{
		UserService: userService,
	}
}

func (handler *LoginHandler) Register(c *fiber.Ctx) error {
	var userRequest request.UserCreateRequest
	err := c.BodyParser(&userRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = helper.ValidateStruct(userRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userCreated, code, err := handler.UserService.Register(c.Context(), &pb.CreateUserRequest{
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: userRequest.Password,
	})
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusCreated, "created", &pb.User{
		Id:        userCreated.GetUser().GetId(),
		Name:      userCreated.GetUser().GetName(),
		Email:     userCreated.GetUser().GetEmail(),
		CreatedAt: userCreated.GetUser().GetCreatedAt(),
		UpdatedAt: userCreated.GetUser().GetUpdatedAt(),
	})
}

func (handler *LoginHandler) Login(c *fiber.Ctx) error {
	var requestLogin request.LoginRequest
	err := c.BodyParser(&requestLogin)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = helper.ValidateStruct(requestLogin)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	signedString, code, err := handler.UserService.ValidateLogin(c.Context(), &pb.GetValidateLoginRequest{
		Email:    requestLogin.Email,
		Password: requestLogin.Password,
	})
	if err != nil {
		return fiber.NewError(code, err.Error())
	}

	return response.ReturnSuccess(c, code, "OK", fiber.Map{
		"access_token": signedString,
	})
}

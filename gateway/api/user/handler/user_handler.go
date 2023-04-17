package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/user/request"
	"log"

	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserController struct {
	UserService pb.UserServiceClient
}

func NewUserController(c *fiber.App, configuration config.Config) *UserController {
	connection, err := grpc.Dial(configuration.Get("USER_SERVICE_URL"), grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}

	controller := &UserController{
		UserService: pb.NewUserServiceClient(connection),
	}

	routes := c.Group("/api")
	routes.Get("/users", controller.FindAll)
	routes.Get("/users/:id", controller.FindByID)
	routes.Post("/users", controller.Create)
	routes.Patch("/users/:id", controller.Update)
	routes.Delete("/users/:id", controller.Delete)

	return controller
}

func (controller *UserController) FindAll(c *fiber.Ctx) error {
	users, err := controller.UserService.FindAll(c.Context(), new(emptypb.Empty))
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessOK(c, "OK", users)
}

func (controller *UserController) FindByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	user, err := controller.UserService.FindByID(c.Context(), &pb.GetUserByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessOK(c, "OK", user)
}

func (controller *UserController) Create(c *fiber.Ctx) error {
	var userRequest request.UserCreateRequest
	err := c.BodyParser(&userRequest)
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	errValidate := helper.ValidateStruct(userRequest)
	if errValidate != nil {
		return response.ReturnErrorBadRequest(c, errValidate)
	}

	userCreated, err := controller.UserService.Create(c.Context(), &pb.CreateUserRequest{
		Name:  userRequest.Name,
		Email: userRequest.Email,
	})
	if err != nil {
		return response.ReturnErrorInternalServerError(c, err)
	}

	return response.ReturnSuccessCreated(c, "created", userCreated)
}

func (controller *UserController) Update(c *fiber.Ctx) error {
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

	userUpdated, err := controller.UserService.Update(c.Context(), &pb.UpdateUserRequest{
		Id:   int64(id),
		Name: userRequest.Name,
	})

	return response.ReturnSuccessOK(c, "updated", userUpdated)
}

func (controller *UserController) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.ReturnErrorBadRequest(c, err)
	}

	_, err = controller.UserService.Delete(c.Context(), &pb.GetUserByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		return err
	}

	return response.ReturnSuccessOK(c, "deleted", nil)
}

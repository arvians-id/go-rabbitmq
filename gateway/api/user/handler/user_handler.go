package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/user/request"
	"github.com/arvians-id/go-rabbitmq/gateway/api/user/services"
	"github.com/goccy/go-json"
	"github.com/rabbitmq/amqp091-go"
	"log"

	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	UserService services.UserServiceContract
	RabbitMQ    *amqp091.Channel
}

func NewUserHandler(userService services.UserServiceContract, rabbitMQ *amqp091.Channel) UserHandler {
	return UserHandler{
		UserService: userService,
		RabbitMQ:    rabbitMQ,
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
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: userRequest.Password,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Send Email To Queue
	_, err = handler.RabbitMQ.QueueDeclare(
		"mail",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		type Message struct {
			ToEmail string
			Message string
		}
		var message Message
		message.ToEmail = userRequest.Email
		message.Message = "Test dynamic message"
		byteMessage, err := json.Marshal(message)
		if err != nil {
			log.Fatalln(err)
		}

		err = handler.RabbitMQ.PublishWithContext(
			c.Context(),
			"",
			"mail",
			false,
			false,
			amqp091.Publishing{
				DeliveryMode: amqp091.Persistent,
				ContentType:  "text/plain",
				Body:         byteMessage,
			},
		)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	return response.ReturnSuccess(c, fiber.StatusCreated, "created", &pb.User{
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

	userUpdated, err := handler.UserService.Update(c.Context(), &pb.UpdateUserRequest{
		Id:       int64(id),
		Name:     userRequest.Name,
		Password: userRequest.Password,
	})
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "updated", &pb.User{
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

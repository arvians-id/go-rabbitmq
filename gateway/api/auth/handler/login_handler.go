package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/auth/request"
	"github.com/arvians-id/go-rabbitmq/gateway/api/auth/services"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type LoginHandler struct {
	UserService services.UserServiceContract
	RabbitMQ    *amqp091.Channel
}

func NewLoginHandler(userService services.UserServiceContract, rabbitMQ *amqp091.Channel) *LoginHandler {
	return &LoginHandler{
		UserService: userService,
		RabbitMQ:    rabbitMQ,
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

func (handler *LoginHandler) Login(c *fiber.Ctx) error {
	var requestLogin request.LoginRequest
	err := c.BodyParser(&requestLogin)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := handler.UserService.ValidateLogin(c.Context(), &pb.GetValidateLoginRequest{
		Email:    requestLogin.Email,
		Password: requestLogin.Password,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	claims := jwt.Claims(jwt.MapClaims{
		"id":    user.GetUser().GetId(),
		"email": user.GetUser().GetEmail(),
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return response.ReturnSuccess(c, 200, "OK", fiber.Map{
		"access_token": signedString,
	})
}

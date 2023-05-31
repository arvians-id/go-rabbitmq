package services

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/client"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

type UserServiceContract interface {
	ValidateLogin(ctx context.Context, in *pb.GetValidateLoginRequest) (string, int, error)
	Register(ctx context.Context, in *pb.CreateUserRequest) (*pb.GetUserResponse, int, error)
	FindAll(ctx context.Context) (*pb.ListUserResponse, int, error)
	FindByIDs(ctx context.Context, in *pb.GetUserByIDsRequest) (*pb.ListUserResponse, int, error)
	FindByID(ctx context.Context, in *pb.GetUserByIDRequest) (*pb.GetUserResponse, int, error)
	Create(ctx context.Context, in *pb.CreateUserRequest) (*pb.GetUserResponse, int, error)
	Update(ctx context.Context, in *pb.UpdateUserRequest) (*pb.GetUserResponse, int, error)
	Delete(ctx context.Context, in *pb.GetUserByIDRequest) (int, error)
}

type userService struct {
	UserClient client.UserClient
	RabbitMQ   *amqp091.Channel
}

func NewUserService(userClient client.UserClient, rabbitMQ *amqp091.Channel) UserServiceContract {
	return &userService{
		UserClient: userClient,
		RabbitMQ:   rabbitMQ,
	}
}
func (service *userService) FindAll(ctx context.Context) (*pb.ListUserResponse, int, error) {
	users, err := service.UserClient.Client.FindAll(ctx, new(emptypb.Empty))
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return users, fiber.StatusOK, nil
}

func (service *userService) FindByIDs(ctx context.Context, in *pb.GetUserByIDsRequest) (*pb.ListUserResponse, int, error) {
	users, err := service.UserClient.Client.FindByIDs(ctx, in)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	return users, fiber.StatusOK, nil
}

func (service *userService) FindByID(ctx context.Context, in *pb.GetUserByIDRequest) (*pb.GetUserResponse, int, error) {
	user, err := service.UserClient.Client.FindByID(ctx, in)
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return nil, fiber.StatusNotFound, err
		}
		return nil, fiber.StatusInternalServerError, err
	}

	return user, fiber.StatusOK, nil
}

func (service *userService) Create(ctx context.Context, in *pb.CreateUserRequest) (*pb.GetUserResponse, int, error) {
	userCreated, err := service.UserClient.Client.Create(ctx, in)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	// Send Email To Queue
	_, err = service.RabbitMQ.QueueDeclare(
		"mail",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	go func() {
		type Message struct {
			ToEmail string
			Message string
		}
		var message Message
		message.ToEmail = userCreated.GetUser().GetEmail()
		message.Message = "Test dynamic message"
		byteMessage, err := json.Marshal(message)
		if err != nil {
			log.Fatalln(err)
		}

		err = service.RabbitMQ.PublishWithContext(
			ctx,
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

	return userCreated, fiber.StatusCreated, nil
}

func (service *userService) Update(ctx context.Context, in *pb.UpdateUserRequest) (*pb.GetUserResponse, int, error) {
	userUpdated, err := service.UserClient.Client.Update(ctx, in)
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return nil, fiber.StatusNotFound, err
		}
		return nil, fiber.StatusInternalServerError, err
	}

	return userUpdated, fiber.StatusOK, nil
}

func (service *userService) Delete(ctx context.Context, in *pb.GetUserByIDRequest) (int, error) {
	_, err := service.UserClient.Client.Delete(ctx, in)
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return fiber.StatusNotFound, err
		}
		return fiber.StatusInternalServerError, err
	}

	return fiber.StatusOK, nil
}

func (service *userService) Register(ctx context.Context, in *pb.CreateUserRequest) (*pb.GetUserResponse, int, error) {
	userCreated, code, err := service.Create(ctx, in)
	if err != nil {
		return nil, code, err
	}

	// Send Email To Queue
	_, err = service.RabbitMQ.QueueDeclare(
		"mail",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fiber.StatusInternalServerError, err
	}

	go func() {
		type Message struct {
			ToEmail string
			Message string
		}
		var message Message
		message.ToEmail = userCreated.GetUser().GetEmail()
		message.Message = "Test dynamic message"
		byteMessage, err := json.Marshal(message)
		if err != nil {
			log.Fatalln(err)
		}

		err = service.RabbitMQ.PublishWithContext(
			ctx,
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

	return userCreated, fiber.StatusCreated, nil
}

func (service *userService) ValidateLogin(ctx context.Context, in *pb.GetValidateLoginRequest) (string, int, error) {
	user, err := service.UserClient.Client.ValidateLogin(ctx, in)
	if err != nil {
		return "", fiber.StatusBadRequest, err
	}

	claims := jwt.Claims(jwt.MapClaims{
		"id":    user.GetUser().GetId(),
		"email": user.GetUser().GetEmail(),
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", fiber.StatusBadRequest, err
	}

	return signedString, fiber.StatusOK, nil
}

package handler

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/dto"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/request"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/services"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

type TodoHandler struct {
	TodoService     services.TodoServiceContract
	CategoryService services.CategoryServiceContract
	RabbitMQ        *amqp091.Channel
}

func NewTodoHandler(todoService services.TodoServiceContract, categoryService services.CategoryServiceContract, rabbitMQ *amqp091.Channel) TodoHandler {
	return TodoHandler{
		TodoService:     todoService,
		CategoryService: categoryService,
		RabbitMQ:        rabbitMQ,
	}
}

func (handler *TodoHandler) DisplayTodoCategoryList(c *fiber.Ctx) error {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var todos *pb.ListTodoResponse
	var categorys *pb.ListCategoryResponse
	var err error
	wg.Add(2)

	go func() {
		var errGo error
		todos, errGo = handler.TodoService.FindAll(c.Context(), new(emptypb.Empty))
		if errGo != nil {
			mutex.Lock()
			err = errGo
			mutex.Unlock()
		}
		defer wg.Done()
	}()

	go func() {
		var errGo error
		categorys, errGo = handler.CategoryService.FindAll(c.Context(), new(emptypb.Empty))
		if errGo != nil {
			mutex.Lock()
			err = errGo
			mutex.Unlock()
		}
		defer wg.Done()
	}()
	wg.Wait()

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "OK", &dto.DisplayCategoryTodoListResponse{
		Todos:      todos.GetTodos(),
		Categories: categorys.GetCategories(),
	})
}

func (handler *TodoHandler) FindAll(c *fiber.Ctx) error {
	todos, err := handler.TodoService.FindAll(c.Context(), new(emptypb.Empty))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "OK", todos.GetTodos())
}

func (handler *TodoHandler) FindByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	todo, err := handler.TodoService.FindByID(c.Context(), &pb.GetTodoByIDRequest{
		Id: int64(id),
	})
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "OK", todo.GetTodo())
}

func (handler *TodoHandler) Create(c *fiber.Ctx) error {
	var todoRequest request.TodoCreateRequest
	err := c.BodyParser(&todoRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = helper.ValidateStruct(todoRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	todoCreated, err := handler.TodoService.Create(c.Context(), &pb.CreateTodoRequest{
		Title:       todoRequest.Title,
		Description: todoRequest.Description,
		UserId:      todoRequest.UserId,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = handler.publish(c.Context(), "todo.created", &dto.DisplayTodoWithCategoriesIDResponse{
		CategoriesID: todoRequest.Categories,
		Id:           todoCreated.GetTodo().GetId(),
		Title:        todoCreated.GetTodo().GetTitle(),
		Description:  todoCreated.GetTodo().GetDescription(),
		IsDone:       proto.Bool(todoCreated.GetTodo().GetIsDone()),
		UserId:       todoCreated.GetTodo().GetUserId(),
		CreatedAt:    todoCreated.GetTodo().GetCreatedAt().AsTime(),
		UpdatedAt:    todoCreated.GetTodo().GetUpdatedAt().AsTime(),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusCreated, "created", todoCreated.GetTodo())
}

func (handler *TodoHandler) Update(c *fiber.Ctx) error {
	var todoRequest request.TodoUpdateRequest
	err := c.BodyParser(&todoRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = helper.ValidateStruct(todoRequest)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	todoUpdated, err := handler.TodoService.Update(c.Context(), &pb.UpdateTodoRequest{
		Id:          int64(id),
		Title:       todoRequest.Title,
		Description: todoRequest.Description,
		IsDone:      todoRequest.IsDone,
		UserId:      todoRequest.UserId,
	})
	if err != nil {
		if err.Error() == response.GrpcErrorNotFound {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return response.ReturnSuccess(c, fiber.StatusOK, "updated", todoUpdated.GetTodo())
}

func (handler *TodoHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	_, err = handler.TodoService.Delete(c.Context(), &pb.GetTodoByIDRequest{
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

func (handler *TodoHandler) publish(ctx context.Context, key string, data interface{}) error {
	marshaled, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = handler.RabbitMQ.PublishWithContext(
		ctx,
		"todo_exchange",
		key,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        marshaled,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

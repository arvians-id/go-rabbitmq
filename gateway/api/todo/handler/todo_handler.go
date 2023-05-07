package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/dto"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/request"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/services"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
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

	return response.ReturnSuccess(c, fiber.StatusOK, "OK", &dto.DisplayCategoryTodoList{
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
	err = handler.RabbitMQ.ExchangeDeclare(
		"todos", // Nama exchange
		"topic", // Jenis exchange
		true,    // Durable
		false,   // Auto-deleted
		false,   // Internal
		false,   // No-wait
		nil,     // Arguments
	)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	type CategoryTodo struct {
		TodoID     int64   `json:"todo_id"`
		CategoryID []int64 `json:"category_id"`
	}
	var categoryTodo CategoryTodo
	categoryTodo.TodoID = todoCreated.GetTodo().GetId()
	categoryTodo.CategoryID = todoRequest.Categories

	data, err := json.Marshal(categoryTodo)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	err = handler.RabbitMQ.PublishWithContext(
		c.Context(),
		"todos",                 // Nama exchange
		"category_todo.created", // Routing key
		false,                   // Mandatory
		false,                   // Immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
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

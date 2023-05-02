package handler

import (
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/dto"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/pb"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/request"
	"github.com/arvians-id/go-rabbitmq/gateway/api/todo/services"
	"github.com/arvians-id/go-rabbitmq/gateway/helper"
	"github.com/arvians-id/go-rabbitmq/gateway/response"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

type TodoHandler struct {
	TodoService         services.TodoServiceContract
	CategoryTodoService services.CategoryTodoServiceContract
}

func NewTodoHandler(todoService services.TodoServiceContract, categoryTodoService services.CategoryTodoServiceContract) TodoHandler {
	return TodoHandler{
		TodoService:         todoService,
		CategoryTodoService: categoryTodoService,
	}
}

func (handler *TodoHandler) DisplayTodoCategoryList(c *fiber.Ctx) error {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var todos *pb.ListTodoResponse
	var categoryTodos *pb.ListCategoryTodoResponse
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
		categoryTodos, errGo = handler.CategoryTodoService.FindAll(c.Context(), new(emptypb.Empty))
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

	// If you want to run this handler sequentially, you can use this code
	//todos, err := handler.TodoService.FindAll(c.Context(), new(emptypb.Empty))
	//if err != nil {
	//	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	//}
	//
	//categoryTodos, err := handler.CategoryTodoService.FindAll(c.Context(), new(emptypb.Empty))
	//if err != nil {
	//	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	//}

	return response.ReturnSuccess(c, fiber.StatusOK, "OK", &dto.DisplayTodoCategoryList{
		Todos:         todos.GetTodos(),
		CategoryTodos: categoryTodos.GetCategoryTodos(),
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
		Title:          todoRequest.Title,
		Description:    todoRequest.Description,
		UserId:         todoRequest.UserId,
		CategoryTodoId: todoRequest.CategoryTodoId,
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

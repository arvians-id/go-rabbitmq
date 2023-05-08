package usecase

import (
	"context"
	"encoding/json"
	"github.com/arvians-id/go-rabbitmq/category_todo/internal/model"
	"github.com/arvians-id/go-rabbitmq/category_todo/internal/repository"
	"github.com/rabbitmq/amqp091-go"
	"sync"
)

type CategoryTodoContract interface {
	Create(ctx context.Context, channel *amqp091.Channel) error
	Delete(ctx context.Context, channel *amqp091.Channel) error
}

type CategoryTodoUsecase struct {
	CategoryTodoRepository repository.CategoryTodoRepositoryContract
}

func NewCategoryTodoUsecase(todoRepository repository.CategoryTodoRepositoryContract) CategoryTodoContract {
	return &CategoryTodoUsecase{
		CategoryTodoRepository: todoRepository,
	}
}
func (usecase *CategoryTodoUsecase) Create(ctx context.Context, channel *amqp091.Channel) error {
	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = channel.QueueBind(
		queue.Name,
		"category_todo.created",
		"category_todo_exchange",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for d := range msgs {
			var categoryTodo model.CategoriesTodo
			err := json.Unmarshal(d.Body, &categoryTodo)
			if err != nil {
				panic(err)
			}
			err = usecase.CategoryTodoRepository.Create(ctx, &model.CategoriesTodo{
				TodoID:     categoryTodo.TodoID,
				CategoryID: categoryTodo.CategoryID,
			})
			if err != nil {
				panic(err)
			}
		}
	}()
	wg.Wait()

	return nil
}

func (usecase *CategoryTodoUsecase) Delete(ctx context.Context, channel *amqp091.Channel) error {
	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = channel.QueueBind(
		queue.Name,
		"category_todo.deleted",
		"category_todo_exchange",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for d := range msgs {
			var categoryTodo model.CategoryTodo
			err := json.Unmarshal(d.Body, &categoryTodo)
			if err != nil {
				panic(err)
			}

			err = usecase.CategoryTodoRepository.Delete(ctx, categoryTodo.TodoID, categoryTodo.CategoryID)
			if err != nil {
				panic(err)
			}
		}
	}()
	wg.Wait()

	return nil
}

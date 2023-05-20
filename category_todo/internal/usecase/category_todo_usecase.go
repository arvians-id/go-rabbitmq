package usecase

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/category_todo/internal/model"
	"github.com/arvians-id/go-rabbitmq/category_todo/internal/repository"
	"github.com/goccy/go-json"
	"github.com/rabbitmq/amqp091-go"
	"sync"
)

type CategoryTodoContract interface {
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

func (usecase *CategoryTodoUsecase) Delete(ctx context.Context, channel *amqp091.Channel) error {
	return usecase.consumeFromExchange(channel, "category_todo_exchange", "category_todo.deleted", func(data []byte) error {
		var categoryTodo model.CategoryTodo
		err := json.Unmarshal(data, &categoryTodo)
		if err != nil {
			return err
		}

		return usecase.CategoryTodoRepository.Delete(ctx, categoryTodo.TodoID, categoryTodo.CategoryID)
	})
}

func (usecase *CategoryTodoUsecase) consumeFromExchange(channel *amqp091.Channel, exchangeName string, routingKey string, consumeFunc func(data []byte) error) error {
	err := channel.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

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
		routingKey,
		exchangeName,
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
	var errs error

	go func() {
		defer wg.Done()
		for d := range msgs {
			err := consumeFunc(d.Body)
			if err != nil {
				errs = err
				return
			}
		}
	}()
	wg.Wait()
	if errs != nil {
		return errs
	}

	return nil
}

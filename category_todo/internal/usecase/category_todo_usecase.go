package usecase

import (
	"context"
	"encoding/json"
	"github.com/arvians-id/go-rabbitmq/category_todo/internal/model"
	"github.com/arvians-id/go-rabbitmq/category_todo/internal/repository"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
)

type CategoryTodoContract interface {
	Create(ctx context.Context, channel *amqp091.Channel, queue amqp091.Queue) (*model.CategoryTodo, error)
	//Delete(ctx context.Context, req *model.CategoryTodo) error
}

type CategoryTodoUsecase struct {
	CategoryTodoRepository repository.CategoryTodoRepository
}

func NewCategoryTodoUsecase(todoRepository repository.CategoryTodoRepository) CategoryTodoContract {
	return &CategoryTodoUsecase{
		CategoryTodoRepository: todoRepository,
	}
}
func (usecase *CategoryTodoUsecase) Create(ctx context.Context, channel *amqp091.Channel, queue amqp091.Queue) (*model.CategoryTodo, error) {
	err := channel.QueueBind(
		queue.Name,              // Nama antrian
		"category_todo.created", // Routing key yang diikat ke antrian
		"todos",                 // Nama exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind a queue: %v", err)
	}

	msgs, err := channel.Consume(
		queue.Name, // Nama antrian
		"",         // Consumer name
		true,       // Auto-acknowledge
		false,      // Exclusive
		false,      // No-local
		false,      // No-wait
		nil,        // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
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
			for _, category := range categoryTodo.CategoryID {
				categoryTodoCreated, err := usecase.CategoryTodoRepository.Create(ctx, &model.CategoryTodoCreate{
					TodoID:     categoryTodo.TodoID,
					CategoryID: category,
				})
				if err != nil {
					panic(err)
				}
				log.Println("Category Todo Created", categoryTodoCreated)
			}
		}
	}()
	wg.Wait()

	return nil, nil
}

//func (usecase *CategoryTodoUsecase) Delete(ctx context.Context, req *model.CategoryTodo) error {
//	err := usecase.CategoryTodoRepository.Delete(ctx, req.TodoID, req.CategoryID)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

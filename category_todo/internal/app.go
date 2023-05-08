package internal

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-rabbitmq/category_todo/internal/repository"
	"github.com/arvians-id/go-rabbitmq/category_todo/internal/usecase"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

func NewApp(ctx context.Context, channel *amqp091.Channel, db *sql.DB) {
	categoryTodoRepository := repository.NewCategoryTodoRepository(db)
	categoryTodoUsecase := usecase.NewCategoryTodoUsecase(categoryTodoRepository)

	go func() {
		err := categoryTodoUsecase.Delete(ctx, channel)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		err := categoryTodoUsecase.Create(ctx, channel)
		if err != nil {
			log.Fatalln(err)
		}
	}()
}

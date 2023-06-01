package internal

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/category_todo/internal/repository"
	"github.com/arvians-id/go-rabbitmq/category_todo/internal/usecase"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"log"
)

func NewApp(ctx context.Context, channel *amqp091.Channel, db *gorm.DB) {
	categoryTodoRepository := repository.NewCategoryTodoRepository(db)
	categoryTodoUsecase := usecase.NewCategoryTodoUsecase(categoryTodoRepository)

	go func() {
		err := categoryTodoUsecase.Delete(ctx, channel)
		if err != nil {
			log.Println(err)
		}
	}()
}

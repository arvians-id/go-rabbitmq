package repository

import (
	"context"
	"gorm.io/gorm"

	"github.com/arvians-id/go-rabbitmq/category_todo/internal/model"
)

type CategoryTodoRepositoryContract interface {
	Delete(ctx context.Context, todoID int64, categoryID int64) error
}

type CategoryTodoRepository struct {
	DB *gorm.DB
}

func NewCategoryTodoRepository(db *gorm.DB) CategoryTodoRepositoryContract {
	return &CategoryTodoRepository{
		DB: db,
	}
}

func (repository *CategoryTodoRepository) Delete(ctx context.Context, todoID int64, categoryID int64) error {
	var categoryTodo model.CategoryTodo
	err := repository.DB.WithContext(ctx).Where("todo_id = $1 AND category_id = $2", todoID, categoryID).Delete(&categoryTodo).Error
	if err != nil {
		return err
	}

	return nil
}

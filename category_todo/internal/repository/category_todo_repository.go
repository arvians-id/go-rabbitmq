package repository

import (
	"context"
	"database/sql"

	"github.com/arvians-id/go-rabbitmq/category_todo/internal/model"
)

type CategoryTodoRepositoryContract interface {
	Create(ctx context.Context, category *model.TodoWithCategoriesIDResponse) error
	Delete(ctx context.Context, todoID int64, categoryID int64) error
	DeleteAllByTodoID(ctx context.Context, todoID int64) error
}

type CategoryTodoRepository struct {
	DB *sql.DB
}

func NewCategoryTodoRepository(db *sql.DB) CategoryTodoRepositoryContract {
	return &CategoryTodoRepository{
		DB: db,
	}
}

func (repository *CategoryTodoRepository) Create(ctx context.Context, category *model.TodoWithCategoriesIDResponse) error {
	for _, categoryID := range category.CategoriesID {
		query := `INSERT INTO category_todo(todo_id, category_id) VALUES($1,$2)`
		_, err := repository.DB.ExecContext(ctx, query, category.Id, categoryID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository *CategoryTodoRepository) Delete(ctx context.Context, todoID int64, categoryID int64) error {
	query := `DELETE FROM category_todo WHERE todo_id = $1 AND category_id = $2`
	_, err := repository.DB.ExecContext(ctx, query, todoID, categoryID)
	if err != nil {
		return err
	}

	return nil
}

func (repository *CategoryTodoRepository) DeleteAllByTodoID(ctx context.Context, todoID int64) error {
	query := `DELETE FROM category_todo WHERE todo_id = $1`
	_, err := repository.DB.ExecContext(ctx, query, todoID)
	if err != nil {
		return err
	}

	return nil
}

package repository

import (
	"context"
	"database/sql"

	"github.com/arvians-id/go-rabbitmq/category_todo/internal/model"
)

type CategoryTodoRepositoryContract interface {
	Create(ctx context.Context, category *model.CategoryTodo) (*model.CategoryTodo, error)
	Delete(ctx context.Context, id int64) error
}

type CategoryTodoRepository struct {
	DB *sql.DB
}

func NewCategoryTodoRepository(db *sql.DB) CategoryTodoRepository {
	return CategoryTodoRepository{
		DB: db,
	}
}

func (repository *CategoryTodoRepository) Create(ctx context.Context, category *model.CategoryTodoCreate) (*model.CategoryTodoCreate, error) {
	query := `INSERT INTO category_todo(todo_id, category_id) VALUES($1,$2)`
	_, err := repository.DB.ExecContext(ctx, query, category.TodoID, category.CategoryID)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (repository *CategoryTodoRepository) Delete(ctx context.Context, todoID int64, categoryID int64) error {
	query := `DELETE FROM category_todo WHERE todo_id = $1 AND category_id = $2`
	_, err := repository.DB.ExecContext(ctx, query, todoID, categoryID)
	if err != nil {
		return err
	}

	return nil
}

package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-rabbitmq/todo/internal/model"
)

type CategoryTodoRepositoryContract interface {
	FindAll(ctx context.Context) ([]*model.CategoryTodo, error)
	FindByID(ctx context.Context, id int64) (*model.CategoryTodo, error)
	Create(ctx context.Context, categoryTodo *model.CategoryTodo) (*model.CategoryTodo, error)
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

func (repository *CategoryTodoRepository) FindAll(ctx context.Context) ([]*model.CategoryTodo, error) {
	query := `SELECT * FROM category_todos ORDER BY created_at DESC`
	rows, err := repository.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categoryTodos []*model.CategoryTodo
	for rows.Next() {
		var categoryTodo model.CategoryTodo
		err := rows.Scan(&categoryTodo.Id, &categoryTodo.Name, &categoryTodo.CreatedAt, &categoryTodo.UpdatedAt)
		if err != nil {
			return nil, err
		}

		categoryTodos = append(categoryTodos, &categoryTodo)
	}

	return categoryTodos, nil
}

func (repository *CategoryTodoRepository) FindByID(ctx context.Context, id int64) (*model.CategoryTodo, error) {
	query := `SELECT * FROM category_todos WHERE id = $1`
	row := repository.DB.QueryRowContext(ctx, query, id)

	var categoryTodo model.CategoryTodo
	err := row.Scan(&categoryTodo.Id, &categoryTodo.Name, &categoryTodo.CreatedAt, &categoryTodo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &categoryTodo, nil
}

func (repository *CategoryTodoRepository) Create(ctx context.Context, categoryTodo *model.CategoryTodo) (*model.CategoryTodo, error) {
	query := `INSERT INTO category_todos(name, created_at, updated_at) VALUES($1,$2,$3) RETURNING id`
	row := repository.DB.QueryRowContext(ctx, query, categoryTodo.Name, categoryTodo.CreatedAt, categoryTodo.UpdatedAt)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	categoryTodo.Id = id

	return categoryTodo, nil
}

func (repository *CategoryTodoRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM category_todos WHERE id = $1`
	_, err := repository.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

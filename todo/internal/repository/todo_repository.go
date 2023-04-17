package repository

import (
	"context"
	"database/sql"

	"github.com/arvians-id/go-rabbitmq/todo/internal/model"
)

type TodoRepositoryContract interface {
	FindAll(ctx context.Context) ([]*model.Todo, error)
	FindByID(ctx context.Context, id int64) (*model.Todo, error)
	Create(ctx context.Context, todo *model.Todo) (*model.Todo, error)
	Update(ctx context.Context, todo *model.Todo) (*model.Todo, error)
	Delete(ctx context.Context, id int64) error
}

type TodoRepository struct {
	DB *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return TodoRepository{
		DB: db,
	}
}

func (repository *TodoRepository) FindAll(ctx context.Context) ([]*model.Todo, error) {
	query := `SELECT * FROM todos ORDER BY created_at DESC`
	rows, err := repository.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*model.Todo
	for rows.Next() {
		var todo model.Todo
		err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.IsDone, &todo.UserId, &todo.CategoryTodoId, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}

		todos = append(todos, &todo)
	}

	return todos, nil
}

func (repository *TodoRepository) FindByID(ctx context.Context, id int64) (*model.Todo, error) {
	query := `SELECT * FROM todos WHERE id = $1`
	row := repository.DB.QueryRowContext(ctx, query, id)

	var todo model.Todo
	err := row.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.IsDone, &todo.UserId, &todo.CategoryTodoId, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (repository *TodoRepository) Create(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	query := `INSERT INTO todos(title, description, user_id, category_todo_id, created_at, updated_at) VALUES($1,$2,$3,$4,$5,$6) RETURNING id`
	row := repository.DB.QueryRowContext(ctx, query, todo.Title, todo.Description, todo.UserId, todo.CategoryTodoId, todo.CreatedAt, todo.UpdatedAt)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	todo.Id = id

	return todo, nil
}

func (repository *TodoRepository) Update(ctx context.Context, todo *model.Todo) (*model.Todo, error) {
	query := `UPDATE todos SET title = $1, description = $2, is_done = $3, user_id = $4, updated_at = $5 WHERE id = $6`
	_, err := repository.DB.ExecContext(ctx, query, todo.Title, todo.Description, todo.IsDone, todo.UserId, todo.UpdatedAt, todo.Id)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (repository *TodoRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM todos WHERE id = $1`
	_, err := repository.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

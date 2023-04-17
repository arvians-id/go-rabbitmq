package repository

import (
	"context"
	"database/sql"

	"github.com/arvians-id/go-rabbitmq/user/internal/model"
)

type UserRepositoryContract interface {
	FindAll(ctx context.Context) ([]*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int64) error
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{
		DB: db,
	}
}

func (repository *UserRepository) FindAll(ctx context.Context) ([]*model.User, error) {
	query := `SELECT * FROM users ORDER BY created_at DESC`
	rows, err := repository.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (repository *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	row := repository.DB.QueryRowContext(ctx, query, id)

	var user model.User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	query := `INSERT INTO users(name, email, created_at, updated_at) VALUES($1,$2,$3,$4) RETURNING id`
	row := repository.DB.QueryRowContext(ctx, query, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	user.Id = id

	return user, nil
}

func (repository *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	query := `UPDATE users SET name = $1, updated_at = $2 WHERE id = $3`
	_, err := repository.DB.ExecContext(ctx, query, user.Name, user.UpdatedAt, user.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := repository.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

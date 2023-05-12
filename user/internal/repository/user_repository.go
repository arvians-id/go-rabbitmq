package repository

import (
	"context"
	"database/sql"

	"github.com/arvians-id/go-rabbitmq/user/cmd/config"
	"github.com/arvians-id/go-rabbitmq/user/internal/model"
	"go.opentelemetry.io/otel"
)

type UserRepositoryContract interface {
	FindAll(ctx context.Context) ([]*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id int64) error
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepositoryContract {
	return &UserRepository{
		DB: db,
	}
}

func (repository *UserRepository) FindAll(ctx context.Context) ([]*model.User, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.UserService/Repository/FindAll")
	defer span.End()

	query := `SELECT * FROM users ORDER BY created_at DESC`
	rows, err := repository.DB.QueryContext(ctxTracer, query)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (repository *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.UserService/Repository/FindByID")
	defer span.End()

	query := `SELECT * FROM users WHERE id = $1`
	row := repository.DB.QueryRowContext(ctxTracer, query, id)

	var user model.User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.UserService/Repository/FindByEmail")
	defer span.End()

	query := `SELECT * FROM users WHERE email = $1`
	row := repository.DB.QueryRowContext(ctxTracer, query, email)

	var user model.User
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.UserService/Repository/Create")
	defer span.End()

	query := `INSERT INTO users(name, email, password, created_at, updated_at) VALUES($1,$2,$3,$4,$5) RETURNING id`
	row := repository.DB.QueryRowContext(ctxTracer, query, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	user.Id = id

	return user, nil
}

func (repository *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.UserService/Repository/Update")
	defer span.End()

	query := `UPDATE users SET name = $1, password = $2, updated_at = $3 WHERE id = $4`
	_, err := repository.DB.ExecContext(ctxTracer, query, user.Name, user.Password, user.UpdatedAt, user.Id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return user, nil
}

func (repository *UserRepository) Delete(ctx context.Context, id int64) error {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.UserService/Repository/Delete")
	defer span.End()

	query := `DELETE FROM users WHERE id = $1`
	_, err := repository.DB.ExecContext(ctxTracer, query, id)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

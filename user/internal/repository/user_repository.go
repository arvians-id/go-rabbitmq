package repository

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/user/cmd/config"
	"github.com/arvians-id/go-rabbitmq/user/internal/model"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type UserRepositoryContract interface {
	FindAll(ctx context.Context) ([]*model.User, error)
	FindByIDs(ctx context.Context, ids []int64) ([]*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryContract {
	return &UserRepository{
		DB: db,
	}
}

func (repository *UserRepository) FindAll(ctx context.Context) ([]*model.User, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.UserService/Repository/FindAll")
	defer span.End()

	var users []*model.User
	query := `SELECT * FROM users ORDER BY created_at DESC`
	err := repository.DB.WithContext(ctxTracer).Raw(query).Scan(&users).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return users, nil
}

func (repository *UserRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.User, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.UserService/Repository/FindByIDs")
	defer span.End()

	var users []*model.User
	query := `SELECT * FROM users WHERE id IN (?) ORDER BY created_at DESC`
	err := repository.DB.WithContext(ctxTracer).Raw(query, ids).Scan(&users).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return users, nil
}

func (repository *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.UserService/Repository/FindByID")
	defer span.End()

	var user model.User
	query := `SELECT id,name,email,created_at,updated_at FROM users WHERE id = ?`
	row := repository.DB.WithContext(ctxTracer).Raw(query, id).Row()
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.UserService/Repository/FindByEmail")
	defer span.End()

	var user model.User
	query := `SELECT * FROM users WHERE email = ?`
	row := repository.DB.WithContext(ctxTracer).Raw(query, email).Row()
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

	err := repository.DB.WithContext(ctxTracer).Create(user).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return user, nil
}

func (repository *UserRepository) Update(ctx context.Context, user *model.User) error {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.UserService/Repository/Update")
	defer span.End()

	err := repository.DB.WithContext(ctxTracer).Updates(user).Error
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

func (repository *UserRepository) Delete(ctx context.Context, id int64) error {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.UserService/Repository/Delete")
	defer span.End()

	var user model.User
	err := repository.DB.WithContext(ctxTracer).Delete(&user, id).Error
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

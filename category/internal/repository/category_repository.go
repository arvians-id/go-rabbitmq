package repository

import (
	"context"
	"database/sql"

	"github.com/arvians-id/go-rabbitmq/category/cmd/config"
	"github.com/arvians-id/go-rabbitmq/category/internal/model"
	"go.opentelemetry.io/otel"
)

type CategoryRepositoryContract interface {
	FindAll(ctx context.Context) ([]*model.Category, error)
	FindByID(ctx context.Context, id int64) (*model.Category, error)
	Create(ctx context.Context, category *model.Category) (*model.Category, error)
	Delete(ctx context.Context, id int64) error
}

type CategoryRepository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return CategoryRepository{
		DB: db,
	}
}

func (repository *CategoryRepository) FindAll(ctx context.Context) ([]*model.Category, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.CategoryService/Repository/FindAll")
	defer span.End()

	query := `SELECT * FROM categories ORDER BY created_at DESC`
	rows, err := repository.DB.QueryContext(ctxTracer, query)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var categorys []*model.Category
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}

		categorys = append(categorys, &category)
	}

	return categorys, nil
}

func (repository *CategoryRepository) FindByID(ctx context.Context, id int64) (*model.Category, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.CategoryService/Repository/FindByID")
	defer span.End()

	query := `SELECT * FROM categories WHERE id = $1`
	row := repository.DB.QueryRowContext(ctxTracer, query, id)

	var category model.Category
	err := row.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &category, nil
}

func (repository *CategoryRepository) Create(ctx context.Context, category *model.Category) (*model.Category, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.CategoryService/Repository/Create")
	defer span.End()

	query := `INSERT INTO categories(name, created_at, updated_at) VALUES($1,$2,$3) RETURNING id`
	row := repository.DB.QueryRowContext(ctxTracer, query, category.Name, category.CreatedAt, category.UpdatedAt)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	category.Id = id

	return category, nil
}

func (repository *CategoryRepository) Delete(ctx context.Context, id int64) error {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.CategoryService/Repository/Delete")
	defer span.End()

	query := `DELETE FROM categories WHERE id = $1`
	_, err := repository.DB.ExecContext(ctxTracer, query, id)
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

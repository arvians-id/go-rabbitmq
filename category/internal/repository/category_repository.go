package repository

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/category/cmd/config"
	"github.com/arvians-id/go-rabbitmq/category/internal/model"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type CategoryRepositoryContract interface {
	FindAll(ctx context.Context) ([]*model.Category, error)
	FindAllByTodoID(ctx context.Context, todoID int64) ([]*model.Category, error)
	FindByIDs(ctx context.Context, ids []int64) ([]*model.Category, error)
	FindByID(ctx context.Context, id int64) (*model.Category, error)
	Create(ctx context.Context, category *model.Category) (*model.Category, error)
	Delete(ctx context.Context, id int64) error
}

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return CategoryRepository{
		DB: db,
	}
}

func (repository *CategoryRepository) FindAll(ctx context.Context) ([]*model.Category, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.CategoryService/Repository/FindAll")
	defer span.End()

	var categories []*model.Category
	query := `SELECT * FROM categories ORDER BY created_at DESC`
	err := repository.DB.WithContext(ctxTracer).Raw(query).Scan(&categories).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return categories, nil
}

func (repository *CategoryRepository) FindAllByTodoID(ctx context.Context, id int64) ([]*model.Category, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.CategoryService/Repository/FindByID")
	defer span.End()

	var category []*model.Category
	query := `SELECT c.* FROM categories c LEFT JOIN category_todo ct ON c.id = ct.category_id WHERE ct.todo_id = ?`
	err := repository.DB.WithContext(ctxTracer).Raw(query, id).Scan(&category).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return category, nil
}

func (repository *CategoryRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.Category, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.CategoryService/Repository/FindByIDs")
	defer span.End()

	var category []*model.Category
	query := `SELECT * FROM categories WHERE id IN (?) ORDER BY created_at DESC`
	err := repository.DB.WithContext(ctxTracer).Raw(query, ids).Scan(&category).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return category, nil
}

func (repository *CategoryRepository) FindByID(ctx context.Context, id int64) (*model.Category, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.CategoryService/Repository/FindByID")
	defer span.End()

	var category model.Category
	query := `SELECT * FROM categories WHERE id = ?`
	row := repository.DB.WithContext(ctxTracer).Raw(query, id).Row()
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

	err := repository.DB.WithContext(ctxTracer).Create(&category).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return category, nil
}

func (repository *CategoryRepository) Delete(ctx context.Context, id int64) error {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.CategoryService/Repository/Delete")
	defer span.End()

	err := repository.DB.WithContext(ctxTracer).Delete(&model.Category{}, id).Error
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

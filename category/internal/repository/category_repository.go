package repository

import (
	"context"
	"gorm.io/gorm"

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
	err := repository.DB.WithContext(ctxTracer).Order("created_at desc").Find(&categories).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return categories, nil
}

func (repository *CategoryRepository) FindByID(ctx context.Context, id int64) (*model.Category, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.CategoryService/Repository/FindByID")
	defer span.End()

	var category model.Category
	err := repository.DB.WithContext(ctxTracer).First(&category, id).Error
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

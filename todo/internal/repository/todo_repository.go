package repository

import (
	"context"
	"github.com/arvians-id/go-rabbitmq/todo/cmd/config"
	"github.com/arvians-id/go-rabbitmq/todo/internal/model"
	"github.com/arvians-id/go-rabbitmq/todo/pb"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type TodoRepositoryContract interface {
	FindAll(ctx context.Context) ([]*model.Todo, error)
	FindByIDs(ctx context.Context, ids []int64) ([]*model.Todo, error)
	FindByID(ctx context.Context, id int64) (*model.Todo, error)
	Create(ctx context.Context, req *pb.CreateTodoRequest) (*model.Todo, error)
	Update(ctx context.Context, req *pb.UpdateTodoRequest) (*model.Todo, error)
	Delete(ctx context.Context, id int64) error
}

type TodoRepository struct {
	DB *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepositoryContract {
	return &TodoRepository{
		DB: db,
	}
}

func (repository *TodoRepository) FindAll(ctx context.Context) ([]*model.Todo, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.TodoService/Repository/FindAll")
	defer span.End()

	var todos []*model.Todo
	query := `SELECT * FROM todos ORDER BY created_at DESC`
	err := repository.DB.WithContext(ctxTracer).Raw(query).Scan(&todos).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return todos, nil
}

func (repository *TodoRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.Todo, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.TodoService/Repository/FindByIDs")
	defer span.End()

	var todos []*model.Todo
	query := `SELECT * FROM todos WHERE id IN (?) ORDER BY created_at DESC`
	err := repository.DB.WithContext(ctxTracer).Raw(query, ids).Scan(&todos).Error
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return todos, nil
}

func (repository *TodoRepository) FindByID(ctx context.Context, id int64) (*model.Todo, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.TodoService/Repository/FindByID")
	defer span.End()

	var todo model.Todo
	query := `SELECT * FROM todos WHERE id = ?`
	row := repository.DB.WithContext(ctxTracer).Raw(query, id).Row()
	err := row.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.IsDone, &todo.UserId, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &todo, nil
}

func (repository *TodoRepository) Create(ctx context.Context, req *pb.CreateTodoRequest) (*model.Todo, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.TodoService/Repository/Create")
	defer span.End()

	var todo model.Todo
	err := repository.DB.WithContext(ctxTracer).Transaction(func(tx *gorm.DB) error {
		var categories []*model.Category
		for _, categoryID := range req.CategoryId {
			var category model.Category
			err := tx.WithContext(ctxTracer).First(&category, categoryID).Error
			if err != nil {
				return err
			}
			categories = append(categories, &model.Category{
				Id: categoryID,
			})
		}

		todo.Categories = categories
		todo.Title = req.Title
		todo.Description = req.Description
		todo.UserId = req.UserId
		err := tx.WithContext(ctxTracer).Select("title", "description", "user_id", "Categories").Create(&todo).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &todo, nil
}

func (repository *TodoRepository) Update(ctx context.Context, req *pb.UpdateTodoRequest) (*model.Todo, error) {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.TodoService/Repository/Update")
	defer span.End()

	var todo model.Todo
	err := repository.DB.WithContext(ctxTracer).Transaction(func(tx *gorm.DB) error {
		var categories []*model.Category
		for _, categoryID := range req.CategoryId {
			var category model.Category
			err := tx.WithContext(ctxTracer).First(&category, categoryID).Error
			if err != nil {
				return err
			}
			categories = append(categories, &model.Category{
				Id: categoryID,
			})
		}

		todo.Id = req.Id
		todo.Categories = categories
		todo.Title = req.Title
		todo.Description = req.Description
		todo.UserId = req.UserId
		todo.IsDone = req.IsDone
		err := tx.WithContext(ctxTracer).Select("title", "description", "is_done", "user_id").Updates(&todo).Error
		if err != nil {
			return err
		}

		err = tx.Model(&todo).Association("Categories").Replace(categories)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &todo, nil
}

func (repository *TodoRepository) Delete(ctx context.Context, id int64) error {
	ctxTracer, span := otel.Tracer(config.ServiceTrace).Start(ctx, "repository.TodoService/Repository/Delete")
	defer span.End()

	err := repository.DB.WithContext(ctxTracer).Delete(&model.Todo{}, id).Error
	if err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

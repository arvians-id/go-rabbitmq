package model

import (
	"time"

	"github.com/arvians-id/go-rabbitmq/user/pb"
)

type Todo struct {
	Id          int64
	Title       string
	Description string
	IsDone      *bool
	UserId      int64
	Categories  []*Category `gorm:"many2many:category_todo;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (todo *Todo) ToPB() *pb.Todo {
	var categories []*pb.Category
	for _, category := range todo.Categories {
		categories = append(categories, category.ToPB())
	}
	todoData := &pb.Todo{
		Id:          todo.Id,
		Title:       todo.Title,
		Description: todo.Description,
		IsDone:      todo.IsDone,
		UserId:      todo.UserId,
		Categories:  categories,
		CreatedAt:   todo.CreatedAt.String(),
		UpdatedAt:   todo.UpdatedAt.String(),
	}
	if todo.CreatedAt.IsZero() && todo.UpdatedAt.IsZero() {
		todoData.CreatedAt = ""
		todoData.UpdatedAt = ""
	}

	return todoData
}

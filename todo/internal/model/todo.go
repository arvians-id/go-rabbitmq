package model

import (
	"time"

	"github.com/arvians-id/go-rabbitmq/todo/pb"
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
	return &pb.Todo{
		Id:          todo.Id,
		Title:       todo.Title,
		Description: todo.Description,
		IsDone:      todo.IsDone,
		UserId:      todo.UserId,
		CreatedAt:   todo.CreatedAt.String(),
		UpdatedAt:   todo.UpdatedAt.String(),
	}
}

package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/arvians-id/go-rabbitmq/todo/pb"
)

type Todo struct {
	Id             int64
	Title          string
	Description    string
	IsDone         *bool
	UserId         int64
	CategoryTodoId int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (todo *Todo) ToPB() *pb.Todo {
	return &pb.Todo{
		Id:             todo.Id,
		Title:          todo.Title,
		Description:    todo.Description,
		IsDone:         todo.IsDone,
		UserId:         todo.UserId,
		CategoryTodoId: todo.CategoryTodoId,
		CreatedAt:      timestamppb.New(todo.CreatedAt),
		UpdatedAt:      timestamppb.New(todo.UpdatedAt),
	}
}

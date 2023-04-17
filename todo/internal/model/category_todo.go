package model

import (
	"github.com/arvians-id/go-rabbitmq/todo/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type CategoryTodo struct {
	Id        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (categoryTodo *CategoryTodo) ToPB() *pb.CategoryTodo {
	return &pb.CategoryTodo{
		Id:        categoryTodo.Id,
		Name:      categoryTodo.Name,
		CreatedAt: timestamppb.New(categoryTodo.CreatedAt),
		UpdatedAt: timestamppb.New(categoryTodo.UpdatedAt),
	}
}

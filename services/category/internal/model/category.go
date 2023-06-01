package model

import (
	"github.com/arvians-id/go-rabbitmq/category/pb"
	"time"
)

type Category struct {
	Id        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (category *Category) ToPB() *pb.Category {
	return &pb.Category{
		Id:        category.Id,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.String(),
		UpdatedAt: category.UpdatedAt.String(),
	}
}

type CategoryWithTodoID struct {
	Id        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	TodoID    int64
}

func (category *CategoryWithTodoID) ToPB() *pb.CategoryWithTodoID {
	return &pb.CategoryWithTodoID{
		Id:        category.Id,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.String(),
		UpdatedAt: category.UpdatedAt.String(),
		TodoId:    category.TodoID,
	}
}

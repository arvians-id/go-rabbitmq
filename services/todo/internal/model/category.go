package model

import (
	"github.com/arvians-id/go-rabbitmq/todo/pb"
	"time"
)

type Category struct {
	Id        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (category *Category) ToPB() *pb.Category {
	categoryData := &pb.Category{
		Id:        category.Id,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.String(),
		UpdatedAt: category.UpdatedAt.String(),
	}

	return categoryData
}

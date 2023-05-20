package model

import (
	"github.com/arvians-id/go-rabbitmq/user/pb"
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
	if category.CreatedAt.IsZero() && category.UpdatedAt.IsZero() {
		categoryData.CreatedAt = ""
		categoryData.UpdatedAt = ""
	}

	return categoryData
}

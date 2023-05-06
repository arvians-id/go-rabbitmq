package model

import (
	"github.com/arvians-id/go-rabbitmq/category/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		CreatedAt: timestamppb.New(category.CreatedAt),
		UpdatedAt: timestamppb.New(category.UpdatedAt),
	}
}

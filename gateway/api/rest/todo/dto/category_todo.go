package dto

import (
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
)

type DisplayCategoryTodoListResponse struct {
	Todos      []*pb.Todo     `json:"todos"`
	Categories []*pb.Category `json:"categories"`
}

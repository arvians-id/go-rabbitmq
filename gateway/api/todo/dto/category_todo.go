package dto

import (
	"github.com/arvians-id/go-rabbitmq/gateway/pb"
	"time"
)

type DisplayCategoryTodoListResponse struct {
	Todos      []*pb.Todo     `json:"todos"`
	Categories []*pb.Category `json:"categories"`
}

type DisplayTodoWithCategoriesIDResponse struct {
	CategoriesID []int64   `json:"categories_id"`
	Id           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	IsDone       *bool     `json:"is_done"`
	UserId       int64     `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

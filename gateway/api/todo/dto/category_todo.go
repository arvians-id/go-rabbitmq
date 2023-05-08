package dto

import "github.com/arvians-id/go-rabbitmq/gateway/pb"

type DisplayCategoryTodoList struct {
	Todos      []*pb.Todo     `json:"todos"`
	Categories []*pb.Category `json:"categories"`
}

type CategoriesTodo struct {
	TodoID     int64   `json:"todo_id"`
	CategoryID []int64 `json:"category_id"`
}
type CategoryTodo struct {
	TodoID     int64 `json:"todo_id"`
	CategoryID int64 `json:"category_id"`
}

package dto

import "github.com/arvians-id/go-rabbitmq/gateway/api/todo/pb"

type DisplayTodoCategoryList struct {
	Todos         []*pb.Todo         `json:"todos"`
	CategoryTodos []*pb.CategoryTodo `json:"category_todos"`
}

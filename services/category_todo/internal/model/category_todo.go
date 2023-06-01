package model

import (
	"time"
)

type CategoryTodo struct {
	TodoID     int64 `json:"todo_id"`
	CategoryID int64 `json:"category_id"`
}

type TodoWithCategoriesIDResponse struct {
	CategoriesID []int64   `json:"categories_id"`
	Id           int64     `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	IsDone       *bool     `json:"is_done"`
	UserId       int64     `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

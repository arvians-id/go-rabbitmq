package model

type CategoryTodo struct {
	TodoID     int64   `json:"todo_id"`
	CategoryID []int64 `json:"category_id"`
}
type CategoryTodoCreate struct {
	TodoID     int64
	CategoryID int64
}

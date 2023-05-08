package model

type CategoriesTodo struct {
	TodoID     int64   `json:"todo_id"`
	CategoryID []int64 `json:"category_id"`
}

type CategoryTodo struct {
	TodoID     int64 `json:"todo_id"`
	CategoryID int64 `json:"category_id"`
}

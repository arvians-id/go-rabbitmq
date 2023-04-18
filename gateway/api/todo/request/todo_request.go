package request

type TodoCreateRequest struct {
	Title          string `json:"title" validate:"required,min=5"`
	Description    string `json:"description" validate:"required"`
	UserId         int64  `json:"user_id" validate:"required,number"`
	CategoryTodoId int64  `json:"category_todo_id" validate:"required,number"`
}

type TodoUpdateRequest struct {
	Title       string `json:"title" validate:"required,min=5"`
	Description string `json:"description" validate:"required"`
	IsDone      *bool  `json:"is_done"`
	UserId      int64  `json:"user_id" validate:"required,number"`
}

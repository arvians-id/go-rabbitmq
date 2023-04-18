package request

type CategoryTodoCreateRequest struct {
	Name string `json:"name" validate:"required,min=3"`
}

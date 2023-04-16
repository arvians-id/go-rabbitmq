package request

type UserCreateRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

type UserUpdateRequest struct {
	Name string `json:"name" validate:"required,min=3"`
}

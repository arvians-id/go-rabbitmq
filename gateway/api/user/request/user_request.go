package request

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserUpdateRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Password string `json:"password" validate:"omitempty,min=6"`
}

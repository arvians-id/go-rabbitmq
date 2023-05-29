// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AuthLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthLoginResponse struct {
	Token string `json:"token"`
}

type AuthRegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthRegisterResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Category struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CategoryCreateRequest struct {
	Name string `json:"name" validate:"required,min=3"`
}

type CategoryFindByIDRequest struct {
	Id int64 `json:"id"`
}

type Todo struct {
	Id          int64       `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	IsDone      bool        `json:"is_done"`
	UserId      int64       `json:"user_id"`
	Categories  []*Category `json:"categories,omitempty"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

type TodoCreateRequest struct {
	Title       string  `json:"title" validate:"required,min=5"`
	Description string  `json:"description" validate:"required"`
	UserId      int64   `json:"user_id" validate:"required,number"`
	Categories  []int64 `json:"categories" validate:"required"`
}

type TodoUpdateRequest struct {
	Title       string  `json:"title" validate:"required,min=5"`
	Description string  `json:"description" validate:"required"`
	IsDone      bool    `json:"is_done"`
	UserId      int64   `json:"user_id" validate:"required,number"`
	Categories  []int64 `json:"categories" validate:"required"`
}

type User struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Password  *string `json:"password,omitempty"`
	Todos     []*Todo `json:"todos,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserUpdateRequest struct {
	Name     string  `json:"name" validate:"required,min=3"`
	Password *string `json:"password,omitempty" validate:"required,min=6"`
}

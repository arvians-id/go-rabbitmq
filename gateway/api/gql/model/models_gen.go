// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginResponse struct {
	Token string `json:"token"`
}

type AuthRegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
	Name string `json:"name"`
}

type CategoryFindByIDRequest struct {
	Id int64 `json:"id"`
}

type DisplayCategoryTodoListResponse struct {
	Todos      []*Todo     `json:"todos"`
	Categories []*Category `json:"categories"`
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
	Title       string  `json:"title"`
	Description string  `json:"description"`
	UserId      int64   `json:"user_id"`
	Categories  []int64 `json:"categories"`
}

type TodoUpdateRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	IsDone      bool    `json:"is_done"`
	UserId      int64   `json:"user_id"`
	Categories  []int64 `json:"categories"`
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
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Name     string  `json:"name"`
	Password *string `json:"password,omitempty"`
}

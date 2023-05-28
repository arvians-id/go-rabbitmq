package model

import (
	"time"

	"github.com/arvians-id/go-rabbitmq/user/pb"
)

type User struct {
	Id        int64
	Name      string
	Email     string
	Password  string
	Todos     []*Todo
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user *User) ToPB() *pb.User {
	var todos []*pb.Todo
	for _, todo := range user.Todos {
		todos = append(todos, todo.ToPB())
	}

	return &pb.User{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
		Todos:     todos,
	}
}

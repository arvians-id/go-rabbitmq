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
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user *User) ToPB() *pb.User {
	return &pb.User{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}

package model

import (
	"time"

	"github.com/arvians-id/go-rabbitmq/user/pb"
)

type User struct {
	Id        int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user *User) ToPB() *pb.User {
	return &pb.User{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		//	CreatedAt: user.ToPB().GetCreatedAt(),
		//	UpdatedAt: user.ToPB().GetUpdatedAt(),
	}
}

package usecase

import (
	"context"
	"log"
	"time"

	"github.com/arvians-id/go-rabbitmq/user/internal/model"
	"github.com/arvians-id/go-rabbitmq/user/internal/repository"
	"github.com/arvians-id/go-rabbitmq/user/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserUsecase struct {
	UserRepository repository.UserRepository
	pb.UnimplementedUserServiceServer
}

func NewUserUsecase(userRepository repository.UserRepository) pb.UserServiceServer {
	return &UserUsecase{
		UserRepository: userRepository,
	}
}

func (usecase *UserUsecase) FindAll(ctx context.Context, empty *emptypb.Empty) (*pb.ListUserResponse, error) {
	users, err := usecase.UserRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var usersPB []*pb.User
	for _, user := range users {
		usersPB = append(usersPB, user.ToPB())
	}

	return &pb.ListUserResponse{
		Users: usersPB,
	}, nil
}

func (usecase *UserUsecase) FindByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserResponse, error) {
	user, err := usecase.UserRepository.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		User: user.ToPB(),
	}, nil
}

func (usecase *UserUsecase) Create(ctx context.Context, req *pb.CreateUserRequest) (*pb.GetUserResponse, error) {
	log.Println(req)
	userCreated, err := usecase.UserRepository.Create(ctx, &model.User{
		Name:      req.GetName(),
		Email:     req.GetEmail(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		User: userCreated.ToPB(),
	}, nil
}

func (usecase *UserUsecase) Update(ctx context.Context, req *pb.UpdateUserRequest) (*pb.GetUserResponse, error) {
	userCheck, err := usecase.UserRepository.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	userUpdated, err := usecase.UserRepository.Update(ctx, &model.User{
		Id:        userCheck.Id,
		Name:      req.GetName(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		User: userUpdated.ToPB(),
	}, nil
}

func (usecase *UserUsecase) Delete(ctx context.Context, req *pb.GetUserByIDRequest) (*emptypb.Empty, error) {
	userCheck, err := usecase.UserRepository.FindByID(ctx, req.Id)
	if err != nil {
		return new(emptypb.Empty), nil
	}

	err = usecase.UserRepository.Delete(ctx, userCheck.Id)
	if err != nil {
		return new(emptypb.Empty), nil
	}

	return new(emptypb.Empty), nil
}

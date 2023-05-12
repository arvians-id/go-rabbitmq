package usecase

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/arvians-id/go-rabbitmq/user/internal/model"
	"github.com/arvians-id/go-rabbitmq/user/internal/repository"
	"github.com/arvians-id/go-rabbitmq/user/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserUsecase struct {
	UserRepository repository.UserRepositoryContract
	pb.UnimplementedUserServiceServer
}

func NewUserUsecase(userRepository repository.UserRepositoryContract) pb.UserServiceServer {
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

func (usecase *UserUsecase) ValidateLogin(ctx context.Context, req *pb.GetValidateLoginRequest) (*pb.GetUserResponse, error) {
	user, err := usecase.UserRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return &pb.GetUserResponse{
		User: user.ToPB(),
	}, nil
}

func (usecase *UserUsecase) Create(ctx context.Context, req *pb.CreateUserRequest) (*pb.GetUserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userCreated, err := usecase.UserRepository.Create(ctx, &model.User{
		Name:      req.GetName(),
		Email:     req.GetEmail(),
		Password:  string(hashedPassword),
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

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		userCheck.Password = string(hashedPassword)
	}

	userUpdated, err := usecase.UserRepository.Update(ctx, &model.User{
		Id:        userCheck.Id,
		Name:      req.GetName(),
		Password:  userCheck.Password,
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
		return new(emptypb.Empty), err
	}

	err = usecase.UserRepository.Delete(ctx, userCheck.Id)
	if err != nil {
		return new(emptypb.Empty), err
	}

	return new(emptypb.Empty), nil
}

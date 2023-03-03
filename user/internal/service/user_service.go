package service

import (
	"context"
	"github.com/arvians-id/go-mongo/user/internal/repository"
	"github.com/arvians-id/go-mongo/user/pb"
	"github.com/arvians-id/go-mongo/user/util"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) pb.UserServiceServer {
	return &UserService{
		UserRepository: userRepository,
	}
}

func (service *UserService) FindAll(ctx context.Context, empty *emptypb.Empty) (*pb.ListResponse, error) {
	users, err := service.UserRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListResponse{
		Users: users,
	}, nil
}

func (service *UserService) FindByID(ctx context.Context, id *pb.GetByIDRequest) (*pb.GetResponse, error) {
	objectID, err := util.ConvertStringToHex(id.ID)
	if err != nil {
		return nil, err
	}

	user, err := service.UserRepository.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return &pb.GetResponse{
		User: user,
	}, nil
}

func (service *UserService) Create(ctx context.Context, request *pb.CreateRequest) (*pb.GetResponse, error) {
	passwordHashed, err := util.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	var user pb.User
	user.ID = util.GenerateID().Hex()
	user.Name = request.Name
	user.Email = request.Email
	user.Password = passwordHashed
	user.CreatedAt = util.PrimitiveDateToTimestampPB()
	user.UpdatedAt = util.PrimitiveDateToTimestampPB()

	userCreated, err := service.UserRepository.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &pb.GetResponse{
		User: userCreated,
	}, nil
}

func (service *UserService) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.GetResponse, error) {
	objectID, err := util.ConvertStringToHex(request.ID)
	if err != nil {
		return nil, err
	}

	checkUser, err := service.UserRepository.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	newPassword := checkUser.Password
	if request.Password != "" {
		passwordHashed, err := util.HashPassword(request.Password)
		if err != nil {
			return nil, err
		}
		newPassword = passwordHashed
	}

	checkUser.ID = request.ID
	checkUser.Name = request.Name
	checkUser.Password = newPassword
	checkUser.UpdatedAt = util.PrimitiveDateToTimestampPB()

	user, err := service.UserRepository.Update(ctx, checkUser)
	if err != nil {
		return nil, err
	}

	return &pb.GetResponse{
		User: user,
	}, nil
}

func (service *UserService) Delete(ctx context.Context, id *pb.GetByIDRequest) (*emptypb.Empty, error) {
	objectID, err := util.ConvertStringToHex(id.ID)
	if err != nil {
		return nil, err
	}

	_, err = service.UserRepository.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	err = service.UserRepository.Delete(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

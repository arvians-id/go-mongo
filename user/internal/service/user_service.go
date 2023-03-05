package service

import (
	"context"
	"github.com/arvians-id/go-mongo/user/internal/model"
	"github.com/arvians-id/go-mongo/user/internal/repository"
	"github.com/arvians-id/go-mongo/user/pb"
	"github.com/arvians-id/go-mongo/util"
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

func (service *UserService) FindAll(ctx context.Context, empty *emptypb.Empty) (*pb.ListUserResponse, error) {
	users, err := service.UserRepository.FindAll(ctx)
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

func (service *UserService) FindByID(ctx context.Context, id *pb.GetUserByIDRequest) (*pb.GetUserResponse, error) {
	objectID, err := util.ConvertStringToHex(id.ID)
	if err != nil {
		return nil, err
	}

	user, err := service.UserRepository.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		User: user.ToPB(),
	}, nil
}

func (service *UserService) Create(ctx context.Context, request *pb.CreateUserRequest) (*pb.GetUserResponse, error) {
	passwordHashed, err := util.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	var user model.User
	user.Name = request.Name
	user.Email = request.Email
	user.Password = passwordHashed
	user.CreatedAt = util.PrimitiveDateTime()
	user.UpdatedAt = util.PrimitiveDateTime()

	userCreated, err := service.UserRepository.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		User: userCreated.ToPB(),
	}, nil
}

func (service *UserService) Update(ctx context.Context, request *pb.UpdateUserRequest) (*pb.GetUserResponse, error) {
	objectID, err := util.ConvertStringToHex(request.ID)
	if err != nil {
		return nil, err
	}

	user, err := service.UserRepository.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	newPassword := user.Password
	if request.Password != "" {
		passwordHashed, err := util.HashPassword(request.Password)
		if err != nil {
			return nil, err
		}
		newPassword = passwordHashed
	}

	user.Name = request.Name
	user.Password = newPassword
	user.UpdatedAt = util.PrimitiveDateTime()

	userUpdated, err := service.UserRepository.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		User: userUpdated.ToPB(),
	}, nil
}

func (service *UserService) Delete(ctx context.Context, id *pb.GetUserByIDRequest) (*emptypb.Empty, error) {
	objectID, err := util.ConvertStringToHex(id.ID)
	if err != nil {
		return nil, err
	}

	user, err := service.UserRepository.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	err = service.UserRepository.Delete(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return new(emptypb.Empty), nil
}

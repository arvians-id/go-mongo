//go:build wireinject
// +build wireinject

package injection

import (
	"github.com/arvians-id/go-mongo/user/cmd/config"
	"github.com/arvians-id/go-mongo/user/internal/repository"
	"github.com/arvians-id/go-mongo/user/internal/service"
	"github.com/arvians-id/go-mongo/user/pb"
	"github.com/google/wire"
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
)

func InitServerAPI(configuration config.Config) (pb.UserServiceServer, error) {
	wire.Build(
		config.NewInitializedDatabase,
		userSet,
	)

	return nil, nil
}

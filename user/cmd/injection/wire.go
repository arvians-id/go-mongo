//go:build wireinject
// +build wireinject

package injection

import (
	"fmt"
	"github.com/arvians-id/go-mongo/user/cmd/config"
	"github.com/arvians-id/go-mongo/user/internal/repository"
	"github.com/arvians-id/go-mongo/user/internal/service"
	"github.com/google/wire"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// Server
var UserSet = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
)

func InitServerAPI(configuration config.Config) (model.UserServiceServer, error) {
	wire.Build(
		config.NewInitializedDatabase,
		UserSet,
	)

	return nil, nil
}

// Service
func ProvidePort(configuration config.Config) (net.Listener, error) {
	port := ":" + strings.Split(configuration.Get("USER_SERVICE_URL"), ":")[1]
	connection, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}
	fmt.Println("User service is running on port", port)

	return connection, nil
}

var ServiceSet = wire.NewSet(
	ConfigSet,
	ProvidePort,
)

func InitService() (net.Listener, error) {
	wire.Build(
		ServiceSet,
	)

	return nil, nil
}

// Configuration
var RootDirSet = wire.NewSet(
	ProvideRootDir,
)

func ProvideRootDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	return wd, nil
}

var ConfigSet = wire.NewSet(
	ProvideConfig,
	RootDirSet,
)

func ProvideConfig(rootDir string) (config.Config, error) {
	configDir := filepath.Join(rootDir, ".env")
	configuration := config.New(configDir)
	return configuration, nil
}

func InitConfig() (config.Config, error) {
	wire.Build(
		ConfigSet,
	)

	return nil, nil
}

// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package injection

import (
	"fmt"
	"github.com/arvians-id/go-mongo/post/cmd/config"
	"github.com/arvians-id/go-mongo/post/internal/repository"
	"github.com/arvians-id/go-mongo/post/internal/service"
	"github.com/arvians-id/go-mongo/post/pb"
	"github.com/google/wire"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// Injectors from wire.go:

func InitServerAPI(configuration config.Config) (pb.PostServiceServer, error) {
	database, err := config.NewInitializedDatabase(configuration)
	if err != nil {
		return nil, err
	}
	postRepository := repository.NewPostRepository(database)
	postServiceServer := service.NewPostService(postRepository)
	return postServiceServer, nil
}

func InitService() (net.Listener, error) {
	string2, err := ProvideRootDir()
	if err != nil {
		return nil, err
	}
	configConfig, err := ProvideConfig(string2)
	if err != nil {
		return nil, err
	}
	listener, err := ProvidePort(configConfig)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func InitConfig() (config.Config, error) {
	string2, err := ProvideRootDir()
	if err != nil {
		return nil, err
	}
	configConfig, err := ProvideConfig(string2)
	if err != nil {
		return nil, err
	}
	return configConfig, nil
}

// wire.go:

// Server
var PostSet = wire.NewSet(repository.NewPostRepository, service.NewPostService)

// Service
func ProvidePort(configuration config.Config) (net.Listener, error) {
	port := ":" + strings.Split(configuration.Get("UserServiceURL"), ":")[1]
	connection, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}
	defer connection.Close()

	return connection, nil
}

var ServiceSet = wire.NewSet(
	ConfigSet,
	ProvidePort,
)

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

package main

import (
	"github.com/arvians-id/go-mongo/user/cmd/config"
	"github.com/arvians-id/go-mongo/user/cmd/injection"
	"github.com/arvians-id/go-mongo/user/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Configuration
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}
	rootDir := filepath.Join(wd, ".env")
	configuration := config.New(rootDir)

	// Port Service
	port := ":" + strings.Split(configuration.Get("UserServiceURL"), ":")[1]
	connection, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("Failed at listening", err)
	}
	defer connection.Close()

	// Server
	server, err := injection.InitServerAPI(configuration)
	if err != nil {
		log.Fatalln("Failed at initializing services", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, server)

	err = grpcServer.Serve(connection)
	if err != nil {
		log.Fatalln("Failed at serving", err)
	}
}

package main

import (
	"github.com/arvians-id/go-mongo/user/cmd/injection"
	"github.com/arvians-id/go-mongo/user/pb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	// Init Configuration
	configuration, err := injection.InitConfig()
	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	// Init Service
	listener, err := injection.InitService()
	if err != nil {
		log.Fatalf("failed to initialize service: %v", err)
	}
	listener.Close()

	// Init Server
	server, err := injection.InitServerAPI(configuration)
	if err != nil {
		log.Fatalln("Failed at initializing services", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, server)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln("Failed at serving", err)
	}
}

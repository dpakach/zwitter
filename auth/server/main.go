package main

import (
	"fmt"
	"log"

	"github.com/dpakach/zwitter/auth/api"
	"github.com/dpakach/zwitter/auth/api/authpb"
	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/pkg/config"
	"github.com/dpakach/zwitter/pkg/service"
	"google.golang.org/grpc"
)

func main() {

	cfg, err := config.NewServerconfig("config/config.yaml")

	if err != nil {
		panic(fmt.Errorf("Failed to read config: %s", err))
	}
	err, UsersEndpoint := cfg.GetNodeAddr("Users")
	if err != nil {
		log.Fatal(err)
	}
	conn, UsersClient := auth.NewUsersClient(UsersEndpoint)
	defer conn.Close()

	service := &service.Service{
		Config:   cfg,
		AuthRPCs: []string{},
		RegisterGrpcServer: func(serv *grpc.Server) {
			authpb.RegisterAuthServiceServer(serv, &api.Server{})
		},
		RegisterRestServer: authpb.RegisterAuthServiceHandlerFromEndpoint,
		UsersServiceClient: UsersClient,
		RPCBasePath:        "/authpb.AuthService/",
	}

	service.Start()
}

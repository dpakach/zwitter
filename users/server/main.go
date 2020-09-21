package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/pkg/config"
	"github.com/dpakach/zwitter/pkg/service"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/dpakach/zwitter/users/api"
	"github.com/dpakach/zwitter/users/api/userspb"
	"google.golang.org/grpc"
)

func main() {

	cfg, err := config.NewServerconfig("config/config.yaml")

	if err != nil {
		fmt.Println(err)
	}


	err, AuthEndpoint := cfg.GetNodeAddr("Auth")
	if err != nil {
		log.Fatal(err)
	}
	conn, AuthClient := auth.NewAuthClient(AuthEndpoint)
	defer conn.Close()

	service := &service.Service{
		Config:   cfg,
		AuthRPCs: []string{},
		RegisterGrpcServer: func(serv *grpc.Server) {
			userspb.RegisterUsersServiceServer(serv, &api.Server{})
		},
		RegisterRestServer: func(ctx context.Context, mux *runtime.ServeMux, grpcAddr string, opts []grpc.DialOption) error {
			return userspb.RegisterUsersServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
		},
		AuthServiceClient: AuthClient,
	}

	service.Start()
}

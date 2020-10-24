package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/pkg/config"
	zlog "github.com/dpakach/zwitter/pkg/log"
	"github.com/dpakach/zwitter/pkg/service"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/dpakach/zwitter/users/api"
	"github.com/dpakach/zwitter/users/api/userspb"
	"google.golang.org/grpc"
)

func main() {

	cfg, err := config.NewServerconfig("config/config.yaml")

	if err != nil {
		panic(fmt.Errorf("Failed to read config: %s", err))
	}

	err, AuthEndpoint := cfg.GetNodeAddr("Auth")
	if err != nil {
		log.Fatal(err)
	}
	conn, AuthClient := auth.NewAuthClient(AuthEndpoint)
	defer conn.Close()

	logger := zlog.New()

	service := &service.Service{
		Config:   cfg,
		AuthRPCs: []string{},
		RegisterGrpcServer: func(serv *grpc.Server, log *zlog.ZwitLogger) {
			userspb.RegisterUsersServiceServer(serv, &api.Server{Log: log})
		},
		RegisterRestServer: func(ctx context.Context, mux *runtime.ServeMux, grpcAddr string, opts []grpc.DialOption) error {
			return userspb.RegisterUsersServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
		},
		AuthServiceClient: AuthClient,
		RPCBasePath:       "/userspb.UsersService/",
		SwaggerFile:       "./swagger/user_api.swagger.json",
		Log:               logger,
	}

	service.Start()
}

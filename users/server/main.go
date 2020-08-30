package main

import (
	"context"
	"fmt"

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

	service := &service.Service{
		Name:        cfg.Server.Name,
		GrpcAddr:    cfg.Server.GrpcAddr,
		RestAddr:    cfg.Server.RestAddr,
		CertFile:    cfg.Server.CertFile,
		KeyFile:     cfg.Server.KeyFile,
		ServerName:  cfg.Server.ServerName,
		RpcBasePath: "/userspb.UsersService/",
		AuthRPCs:    []string{},
		RegisterGrpcServer: func(serv *grpc.Server) {
			userspb.RegisterUsersServiceServer(serv, &api.Server{})
		},
		RegisterRestServer: func(ctx context.Context, mux *runtime.ServeMux, grpcAddr string, opts []grpc.DialOption) error {
			return userspb.RegisterUsersServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
		},
	}

	service.Start()
}

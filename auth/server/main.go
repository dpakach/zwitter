package main

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/dpakach/zwitter/auth/api"
	"github.com/dpakach/zwitter/auth/api/authpb"
	"github.com/dpakach/zwitter/pkg/config"
	"github.com/dpakach/zwitter/pkg/service"
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
		RpcBasePath: "/authpb.AuthService/",
		AuthRPCs:    []string{},
		RegisterGrpcServer: func(serv *grpc.Server) {
			authpb.RegisterAuthServiceServer(serv, &api.Server{})
		},
		RegisterRestServer: func(ctx context.Context, mux *runtime.ServeMux, grpcAddr string, opts []grpc.DialOption) error {
			return authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
		},
	}

	service.Start()
}

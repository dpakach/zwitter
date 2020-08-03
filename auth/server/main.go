package main

import (
  "context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/dpakach/zwitter/pkg/service"
	"github.com/dpakach/zwitter/auth/api"
	"github.com/dpakach/zwitter/auth/api/authpb"
	"google.golang.org/grpc"
)

func main() {
  service := &service.Service{
    Name: "Auth",
    GrpcAddr : fmt.Sprintf("%s:%d", "localhost", 9999),
    RestAddr : fmt.Sprintf("%s:%d", "localhost", 9990),
    CertFile : "cert/server.crt",
    KeyFile : "cert/server.key",
    ServerName: "grpcserver",
    RpcBasePath : "/authpb.AuthService/",
    AuthRPCs: []string{
    },
    RegisterGrpcServer: func(serv *grpc.Server) {
      authpb.RegisterAuthServiceServer(serv, &api.Server{})
    },
    RegisterRestServer: func(ctx context.Context, mux *runtime.ServeMux, grpcAddr string, opts []grpc.DialOption) error {
      return authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
    },
  }

  service.Start()
}

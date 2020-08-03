package main

import (
  "context"
	"fmt"

	"github.com/dpakach/zwitter/pkg/service"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/dpakach/zwitter/users/api"
	"github.com/dpakach/zwitter/users/api/userspb"
	"google.golang.org/grpc"
)

func main() {
  service := &service.Service{
    Name: "Users",
    GrpcAddr : fmt.Sprintf("%s:%d", "localhost", 8888),
    RestAddr : fmt.Sprintf("%s:%d", "localhost", 8889),
    CertFile : "cert/server.crt",
    KeyFile : "cert/server.key",
    ServerName: "grpcserver",
    RpcBasePath : "/userspb.UsersService/",
    AuthRPCs: []string{},
    RegisterGrpcServer: func(serv *grpc.Server) {
      userspb.RegisterUsersServiceServer(serv, &api.Server{})
    },
    RegisterRestServer: func(ctx context.Context, mux *runtime.ServeMux, grpcAddr string, opts []grpc.DialOption) error {
      return userspb.RegisterUsersServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
    },
  }

  service.Start()
}

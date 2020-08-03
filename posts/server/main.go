package main

import (
  "context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/dpakach/zwitter/pkg/service"
	"github.com/dpakach/zwitter/posts/api"
	"github.com/dpakach/zwitter/posts/api/postspb"
	"google.golang.org/grpc"
)

func main() {
  service := &service.Service{
    Name: "Posts",
    GrpcAddr : fmt.Sprintf("%s:%d", "localhost", 7777),
    RestAddr : fmt.Sprintf("%s:%d", "localhost", 7778),
    CertFile : "cert/server.crt",
    KeyFile : "cert/server.key",
    ServerName: "grpcserver",
    RpcBasePath : "/postspb.PostsService/",
    AuthRPCs: []string{
      "CreatePost",
    },
    RegisterGrpcServer: func(serv *grpc.Server) {
      postspb.RegisterPostsServiceServer(serv, &api.Server{})
    },
    RegisterRestServer: func(ctx context.Context, mux *runtime.ServeMux, grpcAddr string, opts []grpc.DialOption) error {
      return postspb.RegisterPostsServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
    },
  }

  service.Start()
}

package main

import (
	"fmt"
	"log"

	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/pkg/config"
	zlog "github.com/dpakach/zwitter/pkg/log"
	"github.com/dpakach/zwitter/pkg/service"
	"github.com/dpakach/zwitter/posts/api"
	"github.com/dpakach/zwitter/posts/api/postspb"
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

	err, AuthEndpoint := cfg.GetNodeAddr("Auth")
	if err != nil {
		log.Fatal(err)
	}
	conn, AuthClient := auth.NewAuthClient(AuthEndpoint)
	defer conn.Close()

	logger := zlog.New()

	service := &service.Service{
		Config: cfg,
		AuthRPCs: []string{
			"CreatePost",
			"LikePost",
		},
		RegisterGrpcServer: func(serv *grpc.Server, log *zlog.ZwitLogger) {
			postspb.RegisterPostsServiceServer(serv, &api.Server{Log: log})
		},
		RegisterRestServer: postspb.RegisterPostsServiceHandlerFromEndpoint,
		UsersServiceClient: UsersClient,
		AuthServiceClient:  AuthClient,
		RPCBasePath:        "/postspb.PostsService/",
		SwaggerFile:        "./swagger/post_api.swagger.json",
		Log:                logger,
	}

	service.Start()
}

package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/metadata"

	"github.com/dpakach/zwitter/auth/api/authpb"
	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/pkg/config"
	"github.com/dpakach/zwitter/posts/api/postspb"
	"github.com/dpakach/zwitter/users/api/userspb"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"
)

type Service struct {
	Config             *config.ServiceConfig
	AuthRPCs           []string
	RPCBasePath        string
	RegisterGrpcServer func(*grpc.Server)
	RegisterRestServer func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
	UsersServiceClient userspb.UsersServiceClient
	PostsServiceClient postspb.PostsServiceClient
	AuthServiceClient  authpb.AuthServiceClient
}

func (s *Service) AuthenticateClient(token string) (*authpb.User, error) {
	if token == "" {
		return nil, fmt.Errorf("Failed to authenticate client: Token not found in the request")
	}
	resp, err := s.AuthServiceClient.AuthenticateToken(context.Background(), &authpb.AuthenticateTokenRequest{Token: token})

	if err != nil {
		return nil, fmt.Errorf("Failed to authenticate client: %v", err)
	}
	log.Printf("authenticated client: %s", resp.User.Username)
	return resp.User, nil
}

func (s *Service) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = context.WithValue(ctx, auth.ConfigKey, s.Config)

	ctx = context.WithValue(ctx, auth.ServiceKey, s)

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		token := strings.Join(md["token"], "")
		for _, method := range s.AuthRPCs {
			if info.FullMethod == s.RPCBasePath+method {
				if token == "" {
					return nil, fmt.Errorf("Failed to authenticate client: Token not found in the request")
				}
				user, err := s.AuthenticateClient(token)
				if err != nil {
					return nil, err
				}
				ctx = context.WithValue(ctx, auth.ClientIDKey, &auth.UserMetaData{Id: user.Id, Username: user.Username})
				return handler(ctx, req)
			}
		}

		// Ignore error for non auth RPCs
		user, _ := s.AuthenticateClient(token)
		if user != nil {
			ctx = context.WithValue(ctx, auth.ClientIDKey, &auth.UserMetaData{Id: user.Id, Username: user.Username})
		}
		return handler(ctx, req)
	}

	return handler(ctx, req)
}

func (s *Service) StartGRPCServer(serv *grpc.Server) error {

	lis, err := net.Listen("tcp", s.Config.Server.GrpcAddr)

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// creds, err := credentials.NewServerTLSFromFile(s.Config.Server.CertFile, s.Config.Server.KeyFile)
	if err != nil {
		return fmt.Errorf("Failed to load TLS keys %v", err)
	}

	opts := []grpc.ServerOption{
		// grpc.Creds(creds),
		grpc.UnaryInterceptor(s.AuthInterceptor),
	}

	serv = grpc.NewServer(opts...)
	s.RegisterGrpcServer(serv)

	log.Printf("starting HTTP/2 gRPC server on %s", s.Config.Server.GrpcAddr)

	if err := serv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %s", err)
	}

	return nil
}

func (s *Service) StartRESTServer(serv *http.Server) error {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(auth.CredMatcher))

	// creds, err := credentials.NewClientTLSFromFile(s.Config.Server.CertFile, s.Config.Server.ServerName)
	// if err != nil {
	// 	return fmt.Errorf("could not load TLS certificate: %s", err)
	// }

	opts := []grpc.DialOption{
		// grpc.WithTransportCredentials((creds))
		grpc.WithInsecure(),
	}
	err := s.RegisterRestServer(ctx, mux, s.Config.Server.GrpcAddr, opts)
	if err != nil {
		return fmt.Errorf("could not register server Ping: %s", err)
	}

	log.Printf("Starting HTTP/1.1 REST server on %s", s.Config.Server.RestAddr)

	serv.Addr = s.Config.Server.RestAddr
	serv.Handler = mux

	err = serv.ListenAndServe()
	if err != nil {
		log.Printf("Could not start REST server")
	}
	return nil
}

func (s *Service) Start() {
	var grpcServer *grpc.Server
	restServer := new(http.Server)
	go func() {
		err := s.StartGRPCServer(grpcServer)
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		err := s.StartRESTServer(restServer)
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	fmt.Println("Recieved terminate, graceful shutdown", sig)

	if restServer != nil {
		tc, timeoutCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer timeoutCancel()
		restServer.Shutdown(tc)
	}
	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
}

package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"

	"github.com/dpakach/zwitter/auth/api/authpb"
	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/pkg/config"
	zlog "github.com/dpakach/zwitter/pkg/log"
	"github.com/dpakach/zwitter/posts/api/postspb"
	"github.com/dpakach/zwitter/users/api/userspb"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"
)

type Service struct {
	Config             *config.ServiceConfig
	AuthRPCs           []string
	RPCBasePath        string
	RegisterGrpcServer func(*grpc.Server, *zlog.ZwitLogger)
	RegisterRestServer func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
	UsersServiceClient userspb.UsersServiceClient
	PostsServiceClient postspb.PostsServiceClient
	AuthServiceClient  authpb.AuthServiceClient
	SwaggerFile        string
	Log                *zlog.ZwitLogger
}

func (s *Service) AuthenticateClient(token string) (*authpb.User, error) {
	if token == "" {
		s.Log.Error("Failed to authenticate a client: Token not found in the response")
		return nil, fmt.Errorf("Failed to authenticate client: Token not found in the request")
	}

	resp, err := s.AuthServiceClient.AuthenticateToken(context.Background(), &authpb.AuthenticateTokenRequest{Token: token})
	if err != nil {
		s.Log.Error("Failed to authenticate a client: Could not contact the auth service")
		return nil, fmt.Errorf("Failed to authenticate client: %v", err)
	}
	if !resp.Auth {
		s.Log.Error("Failed to authenticate a client: Token is not valid")
		return nil, fmt.Errorf("Failed to authenticate client: Token not valid")
	}
	s.Log.Info("Authenticated a new client")
	return resp.User, nil
}

func (s *Service) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = context.WithValue(ctx, auth.ConfigKey, s.Config)
	ctx = context.WithValue(ctx, auth.ServiceKey, s)

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		s.Log.Infof("Recieved gRPC request for %v", info.FullMethod)
		token := strings.Join(md["token"], "")
		for _, method := range s.AuthRPCs {
			if info.FullMethod == s.RPCBasePath+method {
				if token == "" {
					s.Log.Warnf("Token not found in request, failed in %v", info.FullMethod)
					return nil, fmt.Errorf("Failed to authenticate client: Token not found in the request")
				}
				user, err := s.AuthenticateClient(token)
				if err != nil {
					s.Log.Warnf("Failed while authenticating user: %v", err)
					return nil, err
				}
				s.Log.Infof("Authenticated a new client successfully for %v", info.FullMethod)
				ctx = context.WithValue(ctx, auth.ClientIDKey, &auth.UserMetaData{Id: user.Id, Username: user.Username})
				return handler(ctx, req)
			} else {
				if token != "" {
					// Ignore error for non auth RPCs
					user, _ := s.AuthenticateClient(token)
					if user != nil {
						ctx = context.WithValue(ctx, auth.ClientIDKey, &auth.UserMetaData{Id: user.Id, Username: user.Username})
					}
				}
				return handler(ctx, req)
			}
		}
	} else {
		s.Log.Errorf("Error reading request metadata for %v", info.FullMethod)
		return nil, errors.New("Failed to read token from request")
	}

	return handler(ctx, req)
}

func (s *Service) StartGRPCServer(serv *grpc.Server) error {

	lis, err := net.Listen("tcp", s.Config.Server.GrpcAddr)

	if err != nil {
		s.Log.Errorf("Failed to start gRPC server for %v service: %v", s.Config.Server.Name, err)
		return fmt.Errorf("failed to listen: %v", err)
	}

	// creds, err := credentials.NewServerTLSFromFile(s.Config.Server.CertFile, s.Config.Server.KeyFile)
	//if err != nil {
	//return fmt.Errorf("Failed to load TLS keys %v", err)
	//}

	opts := []grpc.ServerOption{
		// grpc.Creds(creds),
		grpc.UnaryInterceptor(s.AuthInterceptor),
	}

	serv = grpc.NewServer(opts...)
	s.RegisterGrpcServer(serv, s.Log)

	s.Log.Infof("starting HTTP/2 gRPC server on %s", s.Config.Server.GrpcAddr)

	if err := serv.Serve(lis); err != nil {
		s.Log.Errorf("Failed while starting gRPC server for %s service on %s: %v", s.Config.Server.Name, s.Config.Server.GrpcAddr, err)
		return fmt.Errorf("failed to serve: %s", err)
	}

	return nil
}

func (s *Service) StartRESTServer(serv *http.Server) error {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/%v/swagger.json", strings.ToLower(s.Config.Server.Name)), func(w http.ResponseWriter, req *http.Request) {
		file, _ := os.Open(s.SwaggerFile)
		io.Copy(w, file)
	})
	gmux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(auth.CredMatcher))

	// creds, err := credentials.NewClientTLSFromFile(s.Config.Server.CertFile, s.Config.Server.ServerName)
	// if err != nil {
	// 	return fmt.Errorf("could not load TLS certificate: %s", err)
	// }

	opts := []grpc.DialOption{
		// grpc.WithTransportCredentials((creds))
		grpc.WithInsecure(),
	}

	s.Log.Infof("Registering REST server for %v service", s.Config.Server.Name)
	err := s.RegisterRestServer(ctx, gmux, s.Config.Server.GrpcAddr, opts)
	if err != nil {
		s.Log.Errorf("Failed registering REST server on gRPC endpoint for %v, %v", s.Config.Server.Name, err)
		return fmt.Errorf("could not register server Ping: %s", err)
	}
	s.Log.Info("Registerd REST server")
	mux.Handle("/", gmux)

	serv.Addr = s.Config.Server.RestAddr
	serv.Handler = mux

	s.Log.Infof("Starting HTTP/1.1 REST server on %s", s.Config.Server.RestAddr)
	err = serv.ListenAndServe()
	if err != nil {
		s.Log.Errorf("Could not start REST server on %v: %v", s.Config.Server.RestAddr, err)
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

	s.Log.Infof("Recieved terminate, graceful shutdown, %s", sig)

	if restServer != nil {
		tc, timeoutCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer timeoutCancel()
		restServer.Shutdown(tc)
	}
	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
}

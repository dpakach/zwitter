package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/dpakach/zwitter/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Service struct {
	Name               string
	GrpcAddr           string
	RestAddr           string
	CertFile           string
	KeyFile            string
	ServerName         string
	AuthRPCs           []string
	RpcBasePath        string
	RegisterGrpcServer func(*grpc.Server)
	RegisterRestServer func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}

func (s *Service) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	for _, method := range s.AuthRPCs {
		if info.FullMethod == s.RpcBasePath+method {
			user, err := auth.AuthenticateClient(ctx)
			if err != nil {
				return nil, err
			}

			ctx = context.WithValue(ctx, auth.ClientIDKey, &auth.UserMetaData{Id: user.Id, Username: user.Username})
			return handler(ctx, req)
		}
	}

	return handler(ctx, req)
}

func (s *Service) StartGRPCServer(serv *grpc.Server) error {

	lis, err := net.Listen("tcp", s.GrpcAddr)

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(s.CertFile, s.KeyFile)
	if err != nil {
		return fmt.Errorf("Failed to load TLS keys %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(creds),
		grpc.UnaryInterceptor(s.AuthInterceptor),
	}

	serv = grpc.NewServer(opts...)
	s.RegisterGrpcServer(serv)

	log.Printf("starting HTTP/2 gRPC server on %s", s.GrpcAddr)

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

	creds, err := credentials.NewClientTLSFromFile(s.CertFile, s.ServerName)
	if err != nil {
		return fmt.Errorf("could not load TLS certificate: %s", err)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials((creds))}
	err = s.RegisterRestServer(ctx, mux, s.GrpcAddr, opts)
	//err = postspb.RegisterPostsServiceHandlerFromEndpoint(ctx, mux, s.GrpcAddr, opts)
	if err != nil {
		return fmt.Errorf("could not register server Ping: %s", err)
	}

	log.Printf("Starting HTTP/1.1 REST server on %s", s.RestAddr)

	serv.Addr = s.RestAddr
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

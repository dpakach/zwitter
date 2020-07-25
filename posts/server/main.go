package main

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

	"github.com/dpakach/zwitter/posts/api"
	"github.com/dpakach/zwitter/posts/api/postspb"

  "github.com/dpakach/zwitter/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func startGRPCServer(address, certFile, keyFile string) (*grpc.Server, error) {

  lis, err := net.Listen("tcp", address)

  if err != nil {
    return nil, fmt.Errorf("failed to listen: %v", err)
  }

  s := api.Server{}

  creds, err := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")
  if err != nil {
    return nil, fmt.Errorf("Failed to load TLS keys %v", err)
  }

  opts := []grpc.ServerOption{
    grpc.Creds(creds),
    grpc.UnaryInterceptor(auth.UnaryInterceptor),
  }

  grpcServer := grpc.NewServer(opts...)
  postspb.RegisterPostsServiceServer(grpcServer, &s)

  log.Printf("starting HTTP/2 gRPC server on %s", address)

  if err := grpcServer.Serve(lis); err != nil {
    return nil, fmt.Errorf("failed to serve: %s", err)
  }

  return grpcServer, nil
}

func startRESTServer(address, grpcAddress, certFile string) (*http.Server, error) {
  ctx := context.Background()

  ctx, cancel := context.WithCancel(ctx)

  defer cancel()

  mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(auth.CredMatcher))

  creds, err := credentials.NewClientTLSFromFile(certFile, "grpcserver")

  if err != nil {
    return nil, fmt.Errorf("could not load TLS certificate: %s", err)
  }

  opts := []grpc.DialOption{grpc.WithTransportCredentials((creds))}

  err = postspb.RegisterPostsServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts)

  if err != nil {
    return nil, fmt.Errorf("could not register server Ping: %s", err)
  }

  log.Printf("Starting HTTP/1.1 REST server on %s", address)

  s := &http.Server{
    Addr: address,
    Handler: mux,
  }

  err = s.ListenAndServe()
  if err != nil {
    log.Printf("Could not start REST server")
  }
  return s, nil
}

func main() {
  grpcAddress := fmt.Sprintf("%s:%d", "localhost", 7777)
  restAddress := fmt.Sprintf("%s:%d", "localhost", 7778)

  certFile := "cert/server.crt"
  keyFile := "cert/server.key"

  grpcServer := make(chan *grpc.Server)
  restServer := make(chan *http.Server)
  go func(s chan *grpc.Server) {
    grpcServer, err  := startGRPCServer(grpcAddress, certFile, keyFile)
    if err != nil {
      log.Fatalf("failed to start gRPC server %s", err)
    }
    s <- grpcServer
  }(grpcServer)

  go func(s chan *http.Server) {
    restServer, err := startRESTServer(restAddress, grpcAddress, certFile)

    fmt.Println(restServer)

    if err != nil {
      log.Fatalf("failed to start REST server %s", err)
    }
    s <- restServer
  }(restServer)

  gs := <-grpcServer
  rs := <-restServer

  sigChan := make(chan os.Signal)
  signal.Notify(sigChan, os.Interrupt)
  signal.Notify(sigChan, os.Kill)

  sig := <-sigChan

  fmt.Println("Recieved terminate, graceful shutdown", sig)

  tc, _ := context.WithTimeout(context.Background(), 30 * time.Second)
  gs.Stop()
  rs.Shutdown(tc)
}

package main

import (
	"context"
	"log"

	"github.com/dpakach/zwitter/users/api/userspb"
  "github.com/dpakach/zwitter/pkg/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
  var conn *grpc.ClientConn

  creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "grpcserver")
  if err != nil {
    log.Fatalf("could not load tls cert: %s", err)
  }

  auth := auth.Authentication{
    Login: "john",
    Password: "doe",
  }

  conn, err = grpc.Dial(":8888", grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&auth))
  if err != nil {
    log.Fatalf("did not connect: %s", err)
  }

  defer conn.Close()

  c := userspb.NewUsersServiceClient(conn)

  response, err := createUser(c, "user2", "newpassword")
  //response, err := getUsers(c)
  //response, err := getUser(c, 1)
  //response, err := authenticate(c, "user2", "newpassword")

  if err != nil {
    log.Fatalf("Erropr when calling SayHello: %s", err)
  }

  log.Printf("Response from server: %s", response)
}

func createUser(c userspb.UsersServiceClient, username, password string) (*userspb.CreateUserResponse, error) {
  return c.CreateUser(context.Background(), &userspb.CreateUserRequest{Username: username, Password: password})
}

func getUsers(c userspb.UsersServiceClient) (*userspb.GetUsersResponse, error) {
  return c.GetUsers(context.Background(), &userspb.EmptyData{})
}

func getUser(c userspb.UsersServiceClient, id int) (*userspb.GetUserResponse, error) {
  return c.GetUser(context.Background(), &userspb.GetUserRequest{Id: int64(id)})
}

func authenticate(c userspb.UsersServiceClient, username, password string) (*userspb.AuthenticateResponse, error) {
  return c.Authenticate(context.Background(), &userspb.AuthenticateRequest{Username: username, Password: password})
}

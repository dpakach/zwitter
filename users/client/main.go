package main

import (
	"context"
	"log"

	"github.com/dpakach/zwitter/users/api/userspb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
  var conn *grpc.ClientConn

  creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "grpcserver")
  if err != nil {
    log.Fatalf("could not load tls cert: %s", err)
  }

  conn, err = grpc.Dial(":8888", grpc.WithTransportCredentials(creds))
  if err != nil {
    log.Fatalf("did not connect: %s", err)
  }

  defer conn.Close()

  c := userspb.NewUsersServiceClient(conn)

  //response, err := createUser(c, "user2", "newpassword")
  //response, err := getUsers(c)
  //response, err := getUser(c, "user2")
  //response, err := getUserById(c, 1)
  response, err := authenticate(c, "user2", "newpassword")

  if err != nil {
    log.Fatalf("Error when calling rpc Server: %s", err)
  }

  log.Printf("Response from server: %s", response)
}

func createUser(c userspb.UsersServiceClient, username, password string) (*userspb.CreateUserResponse, error) {
  return c.CreateUser(context.Background(), &userspb.CreateUserRequest{Username: username, Password: password})
}

func getUsers(c userspb.UsersServiceClient) (*userspb.GetUsersResponse, error) {
  return c.GetUsers(context.Background(), &userspb.EmptyData{})
}

func getUser(c userspb.UsersServiceClient, username string) (*userspb.GetUserResponse, error) {
  return c.GetUser(context.Background(), &userspb.GetUserRequest{Username: username})
}

func getUserById(c userspb.UsersServiceClient, id int) (*userspb.GetUserResponse, error) {
  return c.GetUserByID(context.Background(), &userspb.GetUserByIDRequest{Id: int64(id)})
}

func authenticate(c userspb.UsersServiceClient, username, password string) (*userspb.AuthenticateResponse, error) {
  return c.Authenticate(context.Background(), &userspb.AuthenticateRequest{Username: username, Password: password})
}

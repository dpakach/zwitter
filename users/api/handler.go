package api

import (
  "context"
  "log"
  "time"

  "github.com/dpakach/zwitter/pkg/data"
  "github.com/dpakach/zwitter/pkg/auth"
  "github.com/dpakach/zwitter/users/api/userspb"

  "google.golang.org/grpc/status"
  "google.golang.org/grpc/codes"
)

type Server struct {}

func (s *Server) SayHello(ctx context.Context, in *userspb.PingMessage) (*userspb.PingMessage, error) {
  log.Printf("Received message %s", in.Greeting)
  return &userspb.PingMessage{Greeting: "Hello from Users service"}, nil
}

func userPb(user *data.User) *userspb.User {
    return  &userspb.User{
      Id: user.ID,
      Username: user.Username,
      Created: user.Created,
    }
}

func (s *Server) CreateUser(ctx context.Context, in *userspb.CreateUserRequest) (*userspb.CreateUserResponse, error) {
  check := data.UserStore.GetByUsername(in.GetUsername())
  if check != nil {
    return nil, status.Errorf(codes.InvalidArgument, "Username already taken")
  }
  ts := time.Now().Unix()
  user := &data.User{Username: in.GetUsername(), Created: int64(ts), Password: auth.NewSHA256(in.GetPassword())}

  data.UserStore.AddDbList(user)

  return &userspb.CreateUserResponse{User: userPb(user)  }, nil
}

func (s *Server) GetUsers(ctx context.Context, in *userspb.EmptyData) (*userspb.GetUsersResponse, error) {
  users := data.UserStore.ReadFromDb()
  resp := []*userspb.User{}
  for _, user := range (*users).Users {
    resp = append(resp, userPb(&user))
  }
  return &userspb.GetUsersResponse{
    Users: resp,
  }, nil
}

func (s *Server) GetUser(ctx context.Context, in *userspb.GetUserRequest) (*userspb.GetUserResponse, error) {
  user := data.UserStore.GetByUsername(in.GetUsername())
  if user == nil {
    return nil, status.Errorf(codes.NotFound, "User not found")
  }
  return &userspb.GetUserResponse{User: userPb(user)}, nil
}

func (s *Server) Authenticate(ctx context.Context, in *userspb.AuthenticateRequest) (*userspb.AuthenticateResponse, error) {
  user := data.UserStore.GetByUsername(in.GetUsername())
  if user == nil {
    return nil, status.Errorf(codes.NotFound, "User not found")
  }
  if auth.NewSHA256(in.GetPassword()) == user.Password {
    return &userspb.AuthenticateResponse{Auth: true, User: userPb(user)}, nil
  }
  return &userspb.AuthenticateResponse{Auth: false}, nil
}

func (s *Server) GetUserByID(ctx context.Context, in *userspb.GetUserByIDRequest) (*userspb.GetUserResponse, error) {
  user := data.UserStore.GetByID(in.GetId())
  if user == nil {
    return nil, status.Errorf(codes.NotFound, "User not found")
  }
  return &userspb.GetUserResponse{User: userPb(user)}, nil
}

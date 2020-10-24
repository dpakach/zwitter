package api

import (
	"context"
	"time"

	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/pkg/data"
	zlog "github.com/dpakach/zwitter/pkg/log"
	"github.com/dpakach/zwitter/users/api/userspb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	Log *zlog.ZwitLogger
}

func (s *Server) SayHello(ctx context.Context, in *userspb.PingMessage) (*userspb.PingMessage, error) {
	s.Log.Infof("Received Hello request with: %v", in.GetGreeting())
	return &userspb.PingMessage{Greeting: "Hello from Users service"}, nil
}

func userPb(user *data.User) *userspb.User {
	return &userspb.User{
		Id:       user.ID,
		Username: user.Username,
		Created:  user.Created,
	}
}

func (s *Server) CreateUser(ctx context.Context, in *userspb.CreateUserRequest) (*userspb.CreateUserResponse, error) {
	check := data.UserStore.GetByUsername(in.GetUsername())
	if check != nil {
		s.Log.Errorf("Failed to create new user: User already exists")
		return nil, status.Errorf(codes.InvalidArgument, "Username already taken")
	}
	ts := time.Now().Unix()
	user := &data.User{Username: in.GetUsername(), Created: int64(ts), Password: auth.NewSHA256(in.GetPassword())}

	data.UserStore.AddDbList(user)
	s.Log.Info("Successfully created a new user")

	return &userspb.CreateUserResponse{User: userPb(user)}, nil
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
		s.Log.Errorf("Failed to Fetch the user: User doesn't exists")
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	return &userspb.GetUserResponse{User: userPb(user)}, nil
}

func (s *Server) Authenticate(ctx context.Context, in *userspb.AuthenticateRequest) (*userspb.AuthenticateResponse, error) {
	user := data.UserStore.GetByUsername(in.GetUsername())
	if user == nil {
		s.Log.Errorf("Failed to Authenticate the user: User doesn't exists")
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	if auth.NewSHA256(in.GetPassword()) == user.Password {
		s.Log.Info("Authenticated the user successfully")
		return &userspb.AuthenticateResponse{Auth: true, User: userPb(user)}, nil
	}
	s.Log.Errorf("Failed to Authenticate the user: Password doesn't match")
	return &userspb.AuthenticateResponse{Auth: false}, nil
}

func (s *Server) GetUserByID(ctx context.Context, in *userspb.GetUserByIDRequest) (*userspb.GetUserResponse, error) {
	user := data.UserStore.GetByID(in.GetId())
	if user == nil {
		s.Log.Errorf("Failed to Fetch the user: User doesn't exists")
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	return &userspb.GetUserResponse{User: userPb(user)}, nil
}

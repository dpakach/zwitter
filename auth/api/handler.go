package api

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	"github.com/dpakach/zwitter/auth/api/authpb"
	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/pkg/data"
	"github.com/dpakach/zwitter/users/api/userspb"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, in *authpb.PingMessage) (*authpb.PingMessage, error) {
	log.Printf("Received message %s", in.Greeting)
	return &authpb.PingMessage{Greeting: "Hello from auth service"}, nil
}

func (s *Server) GetToken(ctx context.Context, in *authpb.GetTokenRequest) (*authpb.GetTokenResponse, error) {
	conn, cl := auth.NewUsersClient()
	defer conn.Close()

	resp, err := cl.Authenticate(context.Background(), &userspb.AuthenticateRequest{Username: in.Username, Password: in.Password})
	if err != nil {
		return nil, fmt.Errorf("Failed to authenticate client: %v", err)
	}
	log.Printf("authenticated client: %s", in.Username)

	user := data.User{Username: resp.User.Username, ID: resp.User.Id, Created: resp.User.Created}
	token, err := auth.CreateToken(user)
	if err != nil {
		return nil, fmt.Errorf("Failed to authenticate client: %v", err)
	}
	refreshToken, err := auth.CreateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("Failed to authenticate client: %v", err)
	}

	return &authpb.GetTokenResponse{Token: token, RefreshToken: refreshToken}, nil
}

func (s *Server) AuthenticateToken(ctx context.Context, in *authpb.AuthenticateTokenRequest) (*authpb.AuthenticateTokenResponse, error) {
	user, err := auth.AuthenticateToken(in.Token)
	if err != nil {
		return nil, fmt.Errorf("Failed to authenticate user: %v", err.Error())
	}
	return &authpb.AuthenticateTokenResponse{User: &authpb.User{Username: user.Username, Id: user.ID}}, nil
}

func (s *Server) RefreshToken(ctx context.Context, in *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	token, err := auth.RefreshAccessToken(in.Token)
	if err != nil {
		return nil, fmt.Errorf("Failed to refresh token: %v", err.Error())
	}
	return &authpb.RefreshTokenResponse{Token: token}, nil
}

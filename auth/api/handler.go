package api

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/dpakach/zwitter/auth/api/authpb"
	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/pkg/data"
	zlog "github.com/dpakach/zwitter/pkg/log"
	"github.com/dpakach/zwitter/pkg/service"
	"github.com/dpakach/zwitter/users/api/userspb"
)

type Server struct {
	Log *zlog.ZwitLogger
}

var ConfigParseError = fmt.Errorf("Error while parsing the config")

var UsersServiceNotFoundError = fmt.Errorf("Users service not configured in config")

func (s *Server) SayHello(ctx context.Context, in *authpb.PingMessage) (*authpb.PingMessage, error) {
	s.Log.Info("Received Hello request")
	return &authpb.PingMessage{Greeting: "Hello from auth service"}, nil
}

func (s *Server) GetToken(ctx context.Context, in *authpb.GetTokenRequest) (*authpb.GetTokenResponse, error) {
	svc, ok := ctx.Value(auth.ServiceKey).(*service.Service)
	if !ok {
		s.Log.Warn("Cannot connect to users service, not configured properly")
		return nil, fmt.Errorf("Not configured properly")
	}

	resp, err := svc.UsersServiceClient.Authenticate(context.Background(), &userspb.AuthenticateRequest{Username: in.Username, Password: in.Password})
	if err != nil {
		s.Log.Warnf("Cannot connect to users service, failed to authenticate: %v", err)
		return nil, fmt.Errorf("Failed to authenticate user: %v", err.Error())
	}
	if !resp.Auth || resp.User == nil {
		return nil, fmt.Errorf("Username or password didn't match %v", err.Error())
	}
	s.Log.Info("Authenticated new client")

	user := data.User{Username: resp.User.Username, ID: resp.User.Id, Created: resp.User.Created}
	token, err := auth.CreateToken(user)
	if err != nil {
		s.Log.Errorf("Cannot create new access token: %v", err)
		return nil, fmt.Errorf("Error while creating the token: %v", err)
	}
	refreshToken, err := auth.CreateRefreshToken(user)
	if err != nil {
		s.Log.Errorf("Cannot create new refresh token: %v", err)
		return nil, fmt.Errorf("Error while creating refresh token: %v", err)
	}

	return &authpb.GetTokenResponse{Token: token, RefreshToken: refreshToken, User: &authpb.User{
		Username: user.Username, Id: user.ID, Created: user.Created,
	}}, nil
}

func (s *Server) AuthenticateToken(ctx context.Context, in *authpb.AuthenticateTokenRequest) (*authpb.AuthenticateTokenResponse, error) {
	s.Log.Info("Authenticating the user using token auth")
	user, err := auth.AuthenticateToken(in.Token)
	if err != nil || user == nil {
		s.Log.Errorf("Cannot authenticate the user: %v", err)
		return nil, fmt.Errorf("Failed to authenticate token: %v", err.Error())
	}
	return &authpb.AuthenticateTokenResponse{Auth: true, User: &authpb.User{Username: user.Username, Id: user.ID}}, nil
}

func (s *Server) RefreshToken(ctx context.Context, in *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	s.Log.Info("Refreshing the access token")
	token, err := auth.RefreshAccessToken(in.Token)
	if err != nil {
		s.Log.Errorf("Cannot refresh the token: %v", err)
		return nil, fmt.Errorf("Failed to refresh token: %v", err.Error())
	}
	return &authpb.RefreshTokenResponse{Token: token}, nil
}

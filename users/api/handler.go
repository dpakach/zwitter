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
	userspb.UnimplementedUsersServiceServer
	Log *zlog.ZwitLogger
}

func (s *Server) SayHello(ctx context.Context, in *userspb.PingMessage) (*userspb.PingMessage, error) {
	s.Log.Infof("Received Hello request with: %v", in.GetGreeting())
	return &userspb.PingMessage{Greeting: "Hello from Users service"}, nil
}

func profilePb(profile *data.Profile) *userspb.Profile {
	return &userspb.Profile{
		UserId:      profile.UserId,
		DisplayName: profile.DisplayName,
		Gender:      userspb.Gender(profile.Gender),
		DateOfBirth: profile.Birthday,
	}
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

	userId := data.UserStore.AddDbList(user)
	profile := &data.Profile{UserId: userId, Gender: int(userspb.Gender_NOT_SPECIFIED)}
	data.ProfileStore.AddDbList(profile)

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

func (s *Server) GetUserProfile(ctx context.Context, in *userspb.GetProfileRequest) (*userspb.GetProfileResponse, error) {
	s.Log.Info("get user profile")
	user := data.UserStore.GetByUsername(in.GetUsername())
	if user == nil {
		s.Log.Errorf("Failed to Fetch the user: User doesn't exists")
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	profile := data.ProfileStore.GetByID(user.ID)
	return &userspb.GetProfileResponse{Profile: profilePb(profile), User: userPb(user)}, nil
}

func (s *Server) GetProfile(ctx context.Context, in *userspb.EmptyData) (*userspb.GetProfileResponse, error) {
	s.Log.Info("get profile")
	user, ok := ctx.Value(auth.ClientIDKey).(*auth.UserMetaData)
	if !ok {
		s.Log.Errorf("Failed to create a new post: User not authenticated properly")
		return nil, status.Errorf(codes.Unauthenticated, "Invalid userid")
	}

	dbUser := data.UserStore.GetByID(user.Id)
	if dbUser == nil {
		s.Log.Errorf("Failed to Fetch the user: User doesn't exists")
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	profile := data.ProfileStore.GetByID(dbUser.ID)
	return &userspb.GetProfileResponse{Profile: profilePb(profile), User: userPb(dbUser)}, nil
}

func (s *Server) SetProfile(ctx context.Context, in *userspb.SetProfileRequest) (*userspb.GetProfileResponse, error) {
	ctxUser, ok := ctx.Value(auth.ClientIDKey).(*auth.UserMetaData)
	if !ok {
		s.Log.Errorf("Failed to set the user profile: User not authenticated properly")
		return nil, status.Errorf(codes.Unauthenticated, "Invalid userid")
	}

	user := data.UserStore.GetByID(ctxUser.Id)
	if user == nil {
		s.Log.Errorf("Failed to Fetch the user: User doesn't exists")
		return nil, status.Errorf(codes.NotFound, "User not found")
	}
	profile := data.ProfileStore.GetByID(user.ID)
	if in.GetProfile().GetDisplayName() != "" {
		profile.DisplayName = in.GetProfile().GetDisplayName()
	}
	if in.GetProfile().GetDateOfBirth() != "" {
		profile.Birthday = in.GetProfile().GetDateOfBirth()
	}
	if in.GetProfile().GetGender() != userspb.Gender_NOT_SPECIFIED {
		profile.Gender = int(in.GetProfile().GetGender())
	}
	return &userspb.GetProfileResponse{Profile: profilePb(profile), User: userPb(user)}, nil
}

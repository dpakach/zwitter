package api

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/pkg/data"
	zlog "github.com/dpakach/zwitter/pkg/log"
	"github.com/dpakach/zwitter/pkg/service"
	"github.com/dpakach/zwitter/posts/api/postspb"
	"github.com/dpakach/zwitter/users/api/userspb"
)

type Server struct {
	postspb.UnimplementedPostsServiceServer
	Log *zlog.ZwitLogger
}

var ConfigParseError = fmt.Errorf("Error while parsing the config")

var UsersServiceNotFoundError = fmt.Errorf("Users service not configured in config")

func postPb(ctx context.Context, post *data.Post, user *userspb.User) *postspb.Post {
	userMd, _ := ctx.Value(auth.ClientIDKey).(*auth.UserMetaData)
	childs := []*postspb.Post{}
	svc, ok := ctx.Value(auth.ServiceKey).(*service.Service)
	if len(post.Children) > 0 {
		if !ok {
			return nil
		}

		for _, child := range post.Children {
			resp, err := svc.UsersServiceClient.GetUserByID(context.Background(), &userspb.GetUserByIDRequest{Id: post.Author})
			if err != nil {
				return nil
			}
			childs = append(childs, postPb(ctx, &child, resp.User))
		}
	}

	rezweetPb := &postspb.Post{}

	if post.RezweetId != 0 {
		rezweetPost := data.PostStore.GetByID(post.RezweetId)
		resp, err := svc.UsersServiceClient.GetUserByID(context.Background(), &userspb.GetUserByIDRequest{Id: rezweetPost.Author})
		if err != nil {
			return nil
		}
		rezweetUser := resp.User
		rezweetPb = &postspb.Post{
			Id:       rezweetPost.ID,
			Text:     rezweetPost.Title,
			Created:  rezweetPost.Created,
			Author:   &postspb.User{Id: rezweetUser.Id, Username: rezweetUser.Username},
			Parentid: rezweetPost.ParentId,
			Likes:    data.LikeStore.GetLikesCount(rezweetPost.ID),
			Media:    rezweetPost.Media,
		}
	}

	resp := &postspb.Post{
		Id:       post.ID,
		Text:     post.Title,
		Created:  post.Created,
		Author:   &postspb.User{Id: user.Id, Username: user.Username},
		Parentid: post.ParentId,
		Children: childs,
		Likes:    data.LikeStore.GetLikesCount(post.ID),
		Media:    post.Media,
		Rezweet:  rezweetPb,
	}

	if userMd != nil {
		userLike := data.LikeStore.GetByPostAndAuthor(post.ID, userMd.Id)
		if userLike != nil {
			resp.Liked = true
		}
	}
	return resp
}

func (s *Server) SayHello(ctx context.Context, in *postspb.PingMessage) (*postspb.PingMessage, error) {
	s.Log.Infof("Received message %s", in.Greeting)
	return &postspb.PingMessage{Greeting: "Hello from posts service"}, nil
}

func (s *Server) CreatePost(ctx context.Context, in *postspb.CreatePostRequest) (*postspb.CreatePostResponse, error) {
	ts := time.Now().Unix()

	user, ok := ctx.Value(auth.ClientIDKey).(*auth.UserMetaData)
	if !ok {
		s.Log.Errorf("Failed to create a new post: User not authenticated properly")
		return nil, errors.New("Invalid userid")
	}

	svc, ok := ctx.Value(auth.ServiceKey).(*service.Service)
	if !ok {
		s.Log.Errorf("Failed to create a new post: Service not configured properly")
		return nil, fmt.Errorf("Not configured properly")
	}

	resp, err := svc.UsersServiceClient.GetUserByID(context.Background(), &userspb.GetUserByIDRequest{Id: user.Id})
	if err != nil {
		s.Log.Errorf("Failed to connect to the users service: %v", err)
		return nil, err
	}

	parentid := in.GetParentid()
	if parentid != 0 {
		if parent := data.PostStore.GetByID(parentid); parent == nil {
			s.Log.Errorf("Could not find the parent post: Bad request")
			return nil, fmt.Errorf("Could not find the parent")
		}
	}

	post := &data.Post{Title: in.GetText(), Created: int64(ts), ParentId: parentid, Media: in.GetMedia()}
	post.Author = resp.User.Id

	data.PostStore.AddDbList(post)
	s.Log.Info("Created a new post")

	return &postspb.CreatePostResponse{
		Post: postPb(ctx, post, resp.User),
	}, nil
}

func (s *Server) GetPosts(ctx context.Context, in *postspb.EmptyData) (*postspb.GetPostsResponse, error) {
	svc, ok := ctx.Value(auth.ServiceKey).(*service.Service)
	if !ok {
		s.Log.Errorf("Failed to create a new post: Service not configured properly")
		return nil, fmt.Errorf("Not configured properly")
	}

	posts := data.PostStore.GetPosts()

	result := []*postspb.Post{}

	for _, post := range posts {
		resp, err := svc.UsersServiceClient.GetUserByID(context.Background(), &userspb.GetUserByIDRequest{Id: post.Author})
		if err != nil {
			s.Log.Errorf("Failed to fetch users from users service")
			return nil, errors.New("Failed while retriving users")
		}
		result = append(result, postPb(ctx, &post, resp.User))
	}

	return &postspb.GetPostsResponse{
		Posts: result,
	}, nil
}

func (s *Server) GetPost(ctx context.Context, in *postspb.GetPostRequest) (*postspb.GetPostResponse, error) {
	post := data.PostStore.GetPostWithChilds(in.Id)

	if post == nil {
		s.Log.Errorf("Could not find the requested post")
		return nil, errors.New("Could not find the post")
	}

	svc, ok := ctx.Value(auth.ServiceKey).(*service.Service)
	if !ok {
		s.Log.Errorf("Failed to create a new post: Service not configured properly")
		return nil, fmt.Errorf("Not configured properly")
	}

	resp, err := svc.UsersServiceClient.GetUserByID(context.Background(), &userspb.GetUserByIDRequest{Id: post.Author})
	if err != nil {
		s.Log.Errorf("Failed to fetch users from users service")
		return nil, errors.New("Failed while retriving user")
	}
	return &postspb.GetPostResponse{
		Post: postPb(ctx, post, resp.User),
	}, nil
}

func (s *Server) GetPostChilds(ctx context.Context, in *postspb.GetPostRequest) (*postspb.GetPostsResponse, error) {
	post := data.PostStore.GetByID(in.Id)

	if post == nil {
		s.Log.Errorf("Could not find the requested post")
		return nil, errors.New("Could not find the post")
	}

	childPosts := data.PostStore.GetPostChilds(in.Id)

	svc, ok := ctx.Value(auth.ServiceKey).(*service.Service)
	if !ok {
		s.Log.Errorf("Failed to create a new post: Service not configured properly")
		return nil, fmt.Errorf("Not configured properly")
	}

	result := []*postspb.Post{}
	for _, post := range childPosts {
		resp, err := svc.UsersServiceClient.GetUserByID(context.Background(), &userspb.GetUserByIDRequest{Id: post.Author})
		if err != nil {
			s.Log.Errorf("Failed to fetch users from users service")
			return nil, errors.New("Failed while retriving user")
		}
		result = append(result, postPb(ctx, &post, resp.User))
	}

	return &postspb.GetPostsResponse{
		Posts: result,
	}, nil
}

func (s *Server) LikePost(ctx context.Context, in *postspb.LikeRequest) (*postspb.EmptyData, error) {
	user, ok := ctx.Value(auth.ClientIDKey).(*auth.UserMetaData)
	if !ok {
		s.Log.Errorf("Failed to get user: Invalid userId")
		return nil, errors.New("Invalid userid")
	}

	post := data.PostStore.GetByID(in.Post)

	if post == nil {
		s.Log.Errorf("Could not find the requested post")
		return nil, errors.New("Could not find the post")
	}

	like := data.LikeStore.GetByPostAndAuthor(in.Post, user.Id)
	if like != nil {
		err := data.LikeStore.DeleteLike(in.Post, user.Id)
		if err != nil {
			s.Log.Errorf("Failed to unlike the post: %v", err)
			return nil, err
		}
		return &postspb.EmptyData{}, nil
	}

	data.LikeStore.AddDbList(&data.Like{Post: in.Post, Author: user.Id})

	return &postspb.EmptyData{}, nil
}

func (s *Server) Rezweet(ctx context.Context, in *postspb.RezweetRequest) (*postspb.CreatePostResponse, error) {
	ts := time.Now().Unix()

	user, ok := ctx.Value(auth.ClientIDKey).(*auth.UserMetaData)
	if !ok {
		s.Log.Errorf("Failed to create a new post: User not authenticated properly")
		return nil, errors.New("Invalid userid")
	}

	svc, ok := ctx.Value(auth.ServiceKey).(*service.Service)
	if !ok {
		s.Log.Errorf("Failed to create a new post: Service not configured properly")
		return nil, fmt.Errorf("Not configured properly")
	}

	resp, err := svc.UsersServiceClient.GetUserByID(context.Background(), &userspb.GetUserByIDRequest{Id: user.Id})
	if err != nil {
		s.Log.Errorf("Failed to connect to the users service: %v", err)
		return nil, err
	}

	post := &data.Post{Title: in.GetText(), Created: int64(ts), RezweetId: in.GetRezweetId()}
	post.Author = resp.User.Id

	data.PostStore.AddDbList(post)
	s.Log.Info("Created a new post")

	return &postspb.CreatePostResponse{
		Post: postPb(ctx, post, resp.User),
	}, nil
}

package api

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"

	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/pkg/data"
	"github.com/dpakach/zwitter/pkg/service"
	"github.com/dpakach/zwitter/posts/api/postspb"
	"github.com/dpakach/zwitter/users/api/userspb"
)

type Server struct{}

var ConfigParseError = fmt.Errorf("Error while parsing the config")

var UsersServiceNotFoundError = fmt.Errorf("Users service not configured in config")

func postPb(post *data.Post, user *userspb.User) *postspb.Post {
	return &postspb.Post{
		Id:       post.ID,
		Text:     post.Title,
		Created:  post.Created,
		Author:   &postspb.User{Id: user.Id, Username: user.Username},
		Parentid: post.ParentId,
	}
}

func (s *Server) SayHello(ctx context.Context, in *postspb.PingMessage) (*postspb.PingMessage, error) {
	log.Printf("Received message %s", in.Greeting)
	return &postspb.PingMessage{Greeting: "Hello from posts service"}, nil
}

func (s *Server) CreatePost(ctx context.Context, in *postspb.CreatePostRequest) (*postspb.CreatePostResponse, error) {
	ts := time.Now().Unix()

	user, ok := ctx.Value(auth.ClientIDKey).(*auth.UserMetaData)
	if !ok {
		return nil, errors.New("Invalid userid")
	}

	svc, ok := ctx.Value(auth.ServiceKey).(*service.Service)
	if !ok {
		return nil, fmt.Errorf("Not configured properly")
	}

	resp, err := svc.UsersServiceClient.GetUserByID(context.Background(), &userspb.GetUserByIDRequest{Id: user.Id})
	if err != nil {
		return nil, err
	}

	parentid := in.GetParentid()
	if parentid != 0 {
		if parent := data.PostStore.GetByID(parentid); parent == nil {
			return nil, fmt.Errorf("Could not find the parent")
		}
	}

	post := &data.Post{Title: in.GetText(), Created: int64(ts), ParentId: parentid}
	post.Author = resp.User.Id

	data.PostStore.AddDbList(post)

	return &postspb.CreatePostResponse{
		Post: postPb(post, resp.User),
	}, nil
}

func (s *Server) GetPosts(ctx context.Context, in *postspb.EmptyData) (*postspb.GetPostsResponse, error) {
	svc, ok := ctx.Value(auth.ServiceKey).(*service.Service)
	if !ok {
		return nil, fmt.Errorf("Not configured properly")
	}

	posts := data.PostStore

	result := []*postspb.Post{}

	for _, post := range posts.Posts {
		resp, err := svc.UsersServiceClient.GetUserByID(context.Background(), &userspb.GetUserByIDRequest{Id: post.Author})
		if err != nil {
			return nil, errors.New("Failed while retriving users")
		}
		result = append(result, postPb(&post, resp.User))
	}

	return &postspb.GetPostsResponse{
		Posts: result,
	}, nil
}

func (s *Server) GetPost(ctx context.Context, in *postspb.GetPostRequest) (*postspb.GetPostResponse, error) {
	post := data.PostStore.GetByID(in.Id)

	if post == nil {
		return nil, errors.New("Could not find the post")
	}

	svc, ok := ctx.Value(auth.ServiceKey).(*service.Service)
	if !ok {
		return nil, fmt.Errorf("Not configured properly")
	}

	resp, err := svc.UsersServiceClient.GetUserByID(context.Background(), &userspb.GetUserByIDRequest{Id: post.Author})
	if err != nil {
		return nil, errors.New("Failed while retriving user")
	}
	return &postspb.GetPostResponse{
		Post: postPb(post, resp.User),
	}, nil
}

func (s *Server) GetPostChilds(ctx context.Context, in *postspb.GetPostRequest) (*postspb.GetPostsResponse, error) {
	post := data.PostStore.GetByID(in.Id)

	if post == nil {
		return nil, errors.New("Could not find the post")
	}

	childPosts := data.PostStore.GetPostChilds(in.Id)

	svc, ok := ctx.Value(auth.ServiceKey).(*service.Service)
	if !ok {
		return nil, fmt.Errorf("Not configured properly")
	}

	result := []*postspb.Post{}
	for _, post := range childPosts {
		resp, err := svc.UsersServiceClient.GetUserByID(context.Background(), &userspb.GetUserByIDRequest{Id: post.Author})
		if err != nil {
			return nil, errors.New("Failed while retriving user")
		}
		result = append(result, postPb(&post, resp.User))
	}

	return &postspb.GetPostsResponse{
		Posts: result,
	}, nil
}

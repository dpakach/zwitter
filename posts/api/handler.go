package api

import (
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"

	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/pkg/data"
	"github.com/dpakach/zwitter/posts/api/postspb"
	"github.com/dpakach/zwitter/users/api/userspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {}

func postPb(post *data.Post, user *userspb.User) *postspb.Post {
    return  &postspb.Post{
      Id: post.ID,
      Text: post.Title,
      Created: post.Created,
      Author: &postspb.User{Id: user.Id, Username: user.Username},
    }
}


func NewUsersClient() (*grpc.ClientConn, userspb.UsersServiceClient) {
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

  c := userspb.NewUsersServiceClient(conn)

  return conn, c
}

func (s *Server) SayHello(ctx context.Context, in *postspb.PingMessage) (*postspb.PingMessage, error) {
  log.Printf("Received message %s", in.Greeting)
  return &postspb.PingMessage{Greeting: "Hello from posts service"}, nil
}

func (s *Server) CreatePost(ctx context.Context, in *postspb.CreatePostRequest) (*postspb.CreatePostResponse, error) {
  ts := time.Now().Unix()
	userid, ok := int64(1), true //ctx.Value("userid").(int64)
	if !ok {
		return nil, errors.New("Invalid userid")
	}

  conn, cl := NewUsersClient()
  defer conn.Close()

  resp, err := cl.GetUser(context.Background(), &userspb.GetUserRequest{Id: userid})
	if err != nil {
    fmt.Println(err)
		return nil, err
	}

	post := &data.Post{Title: in.GetText(), Created: int64(ts)}
	post.Author = resp.User.Id

  data.PostStore.AddDbList(post)

	return &postspb.CreatePostResponse{
    Post: postPb(post, resp.User),
	}, nil
}

func (s *Server) GetPosts(ctx context.Context, in *postspb.EmptyData) (*postspb.GetPostsResponse, error) {
  conn, cl := NewUsersClient()
  defer conn.Close()

  posts := data.PostStore

  result := []*postspb.Post{}

  for _, post := range posts.Posts{
    resp, err := cl.GetUser(context.Background(), &userspb.GetUserRequest{Id: post.Author})
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

  conn, cl := NewUsersClient()
  defer conn.Close()

  resp, err := cl.GetUser(context.Background(), &userspb.GetUserRequest{Id: post.Author})
  if err != nil {
    return nil, errors.New("Failed while retriving user")
  }
	return &postspb.GetPostResponse{
    Post: postPb(post, resp.User),
	}, nil
}

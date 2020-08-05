package main

import (
	"context"
	"log"

	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/posts/api/postspb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	var conn *grpc.ClientConn

	creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "grpcserver")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	// auth := auth.BasicAuthentication{
	//   Login: "user2",
	//   Password: "newpassword",
	// }

	authbearer := &auth.BearerAuthentication{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE1OTY0MjYzMDQsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoidXNlcjIifQ.QWhspHFiGEDCONk_ocO5AV9Td3xXbAG_1VTGyI8iSFg",
	}

	conn, err = grpc.Dial(":7777", grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(authbearer))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	defer conn.Close()

	c := postspb.NewPostsServiceClient(conn)

	//response, err := c.SayHello(context.Background(), &postspb.PingMessage{Greeting: "foo"})
	response, err := createPost(c, "This is a Hello Post")
	//response, err := getPosts(c)
	//response, err := getPost(c, 1)

	if err != nil {
		log.Fatalf("Error when calling rpc method: %s", err)
	}

	log.Printf("Response from server: %s", response)
}

func createPost(c postspb.PostsServiceClient, text string) (*postspb.CreatePostResponse, error) {
	return c.CreatePost(context.Background(), &postspb.CreatePostRequest{Text: text})
}

func getPosts(c postspb.PostsServiceClient) (*postspb.GetPostsResponse, error) {
	return c.GetPosts(context.Background(), &postspb.EmptyData{})
}

func getPost(c postspb.PostsServiceClient, id int) (*postspb.GetPostResponse, error) {
	return c.GetPost(context.Background(), &postspb.GetPostRequest{Id: int64(id)})
}

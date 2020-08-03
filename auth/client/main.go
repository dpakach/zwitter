package main

import (
	"context"
	"log"

	"github.com/dpakach/zwitter/auth/api/authpb"
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

	conn, err = grpc.Dial(":9999", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	defer conn.Close()

	c := authpb.NewAuthServiceClient(conn)

	//response, err := c.SayHello(context.Background(), &authpb.PingMessage{Greeting: "foo"})
	response, err := c.GetToken(context.Background(), &authpb.GetTokenRequest{Username: "user2", Password: "newpassword"})
	//response, err := c.AuthenticateToken(context.Background(), &authpb.AuthenticateTokenRequest{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InVzZXIyIn0.BFJDkjokehcgZv59OFs4ztk8lW2GFYZa4cg_9ddInRo"})
	//response, err := c.RefreshToken(context.Background(), &authpb.RefreshTokenRequest{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTYzNzkxNzcsInVzZXJfaWQiOjEsInVzZXJuYW1lIjoidXNlcjIiLCJ1dWlkIjoiMTBiMTBkNDItYjk5MS00MGMwLTdkYWMtNjFiYmI4YjQxNzY2In0.Hy2YnAEzpC0xiE-nGeuHR-8KBukvTdAi1dnxVQ_R0Cg"})

	if err != nil {
		log.Fatalf("Error when calling rpc method: %s", err)
	}

	log.Printf("Response from server: %s", response)
}

func createPost(c postspb.PostsServiceClient, text string) (*postspb.CreatePostResponse, error) {
	return c.CreatePost(context.Background(), &postspb.CreatePostRequest{Text: text})
}

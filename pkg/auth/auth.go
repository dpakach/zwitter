package auth

import (
	"context"
	"fmt"
	"log"
	"strings"
  "crypto/sha256"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/dpakach/zwitter/users/api/userspb"

	"google.golang.org/grpc/credentials"
)

type contextKey int

const (
  ClientIDKey contextKey = iota
)

type UserMetaData struct {
  Id int64
  Username string
}

type Authentication struct {
  Login string
  Password string
}

// NewSHA256 ...
func NewSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x\n", hash)
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
  return map[string]string {
    "login": a.Login,
    "password": a.Password,
  }, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
  return true
}

func CredMatcher(headerName string) (mdName string, ok bool) {
  if headerName == "Login" || headerName == "Password" {
    return headerName, true
  }
  return "", false
}

func authenticateClient(ctx context.Context) (*userspb.User, error) {
  if md, ok := metadata.FromIncomingContext(ctx); ok {
    clientLogin := strings.Join(md["login"], "")
    clientPassword := strings.Join(md["password"], "")

    conn, cl := NewUsersClient()
    defer conn.Close()

    resp, err := cl.Authenticate(context.Background(), &userspb.AuthenticateRequest{Username: clientLogin, Password: clientPassword})

    if err != nil {
      return nil, fmt.Errorf("Failed to authenticate client: %v", err)
    }

    log.Printf("authenticated client: %s", clientLogin)
    return resp.User, nil
  }

  return nil, fmt.Errorf("missing credentials")
}

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
  user, err := authenticateClient(ctx)
  if err != nil {
    return nil, err
  }

  ctx = context.WithValue(ctx, ClientIDKey, &UserMetaData{Id: user.Id, Username: user.Username})
  return handler(ctx, req)
}

func NewUsersClient() (*grpc.ClientConn, userspb.UsersServiceClient) {
  var conn *grpc.ClientConn

  creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "grpcserver")
  if err != nil {
    log.Fatalf("could not load tls cert: %s", err)
  }

  conn, err = grpc.Dial(":8888", grpc.WithTransportCredentials(creds))
  if err != nil {
    log.Fatalf("did not connect: %s", err)
  }

  c := userspb.NewUsersServiceClient(conn)

  return conn, c
}

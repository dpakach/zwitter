package auth

import (
	"context"
	"fmt"
	"log"
	"strings"
  "crypto/sha256"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey int

const (
  clientIDKey contextKey = iota
)

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

func authenticateClient(ctx context.Context) (string, error) {
  if md, ok := metadata.FromIncomingContext(ctx); ok {
    clientLogin := strings.Join(md["login"], "")
    clientPassword := strings.Join(md["password"], "")

    if clientLogin != "john" {
      return "", fmt.Errorf("unknown user %s", clientLogin)
    }

    if clientPassword != "doe" {
      return "", fmt.Errorf("bad Password %s", clientPassword)
    }

    log.Printf("authenticated client: %s", clientLogin)
    return "42", nil
  }

  return "", fmt.Errorf("missing credentials")
}

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
  clientID, err := authenticateClient(ctx)
  if err != nil {
    return nil, err
  }

  ctx = context.WithValue(ctx, clientIDKey, clientID)
  return handler(ctx, req)
}


package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/dgrijalva/jwt-go"
	"github.com/dpakach/zwitter/auth/api/authpb"
	"github.com/dpakach/zwitter/pkg/data"
	"github.com/dpakach/zwitter/users/api/userspb"

	"google.golang.org/grpc/credentials"
)

var refreshTokens = map[string]string{}

type contextKey int

const (
	ClientIDKey contextKey = iota
)

type UserMetaData struct {
	Id       int64
	Username string
}

type BearerAuthentication struct {
	Token string
}

func (a *BearerAuthentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"token": a.Token,
	}, nil
}

func (a *BearerAuthentication) RequireTransportSecurity() bool {
	return true
}

type BasicAuthentication struct {
	Login    string
	Password string
}

func (a *BasicAuthentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login":    a.Login,
		"password": a.Password,
	}, nil
}

func (a *BasicAuthentication) RequireTransportSecurity() bool {
	return true
}

// NewSHA256 ...
func NewSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x\n", hash)
}

func UUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("could not generate UUID")
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}

func CredMatcher(headerName string) (mdName string, ok bool) {
	if headerName == "Login" || headerName == "Password" {
		return headerName, true
	}
	return "", false
}

func AuthenticateClient(ctx context.Context) (*authpb.User, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		token := strings.Join(md["token"], "")
		if token == "" {
			return nil, fmt.Errorf("Failed to authenticate client: Token not found in the request")
		}
		conn, cl := NewAuthClient()
		defer conn.Close()

		resp, err := cl.AuthenticateToken(context.Background(), &authpb.AuthenticateTokenRequest{Token: token})

		if err != nil {
			return nil, fmt.Errorf("Failed to authenticate client: %v", err)
		}

		log.Printf("authenticated client: %s", resp.User.Username)
		return resp.User, nil
	}

	return nil, fmt.Errorf("missing credentials")
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

func NewAuthClient() (*grpc.ClientConn, authpb.AuthServiceClient) {
	var conn *grpc.ClientConn

	creds, err := credentials.NewClientTLSFromFile("cert/server.crt", "grpcserver")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	conn, err = grpc.Dial(":9999", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	c := authpb.NewAuthServiceClient(conn)

	return conn, c
}

func CreateRefreshToken(user data.User) (string, error) {
	// TODO fix env var
	os.Setenv("REFRESH_SECRET", "my_another_super_secret")

	tokenUuid, _ := UUID()

	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = user.ID
	atClaims["username"] = user.Username
	atClaims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	atClaims["uuid"] = tokenUuid

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return "", err
	}

	// Saving token temporary
	// Need database or cache of some sort for this
	refreshTokens[tokenUuid] = token

	return token, nil
}

func AuthenticateToken(accessToken string) (*data.User, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	user, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Failed to authenticate client: %v", err)
	}
	auth, auth_ok := user["authorized"].(bool)
	username, uname_ok := user["username"].(string)
	user_id, uid_ok := user["user_id"].(float64)
	if !auth_ok || !uid_ok || !uname_ok {
		return nil, fmt.Errorf("Failed to authenticate client: Invalid Token")
	}
	if auth {
		return &data.User{Username: username, ID: int64(user_id)}, nil
	}
	return nil, fmt.Errorf("Failed to authenticate client")
}

func CreateToken(user data.User) (string, error) {
	// TODO fix env var
	os.Setenv("ACCESS_SECRET", "my_super_secret")
	os.Setenv("REFRESH_SECRET", "my_another_super_secret")

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = user.ID
	atClaims["username"] = user.Username
	atClaims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func RefreshAccessToken(refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil {
		return "", fmt.Errorf("Refresh token expired: %v", err)
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return "", fmt.Errorf("Refresh token invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uuid, ok := claims["uuid"].(string)
		username, uname_ok := claims["username"].(string)
		user_id, uid_ok := claims["user_id"].(float64)
		if !ok || !uname_ok || !uid_ok {
			return "", fmt.Errorf("Refresh token invalid")
		}

		delete(refreshTokens, uuid)
		accessToken, err := CreateToken(data.User{Username: username, ID: int64(user_id)})
		if err != nil {
			return "", err
		}

		return accessToken, nil
	}

	return "", fmt.Errorf("Refresh token not valid")
}

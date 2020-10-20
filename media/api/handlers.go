package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dpakach/zwitter/auth/api/authpb"
	"github.com/dpakach/zwitter/media/storage"
	"github.com/dpakach/zwitter/pkg/config"
	"github.com/dpakach/zwitter/pkg/data"
	zlog "github.com/dpakach/zwitter/pkg/log"
	"github.com/gorilla/mux"
)

type KeyUserData struct{}

type Hello struct {
	Log *zlog.ZwitLogger
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(200)
	h.Log.Info("Received hello request")
	rw.Write([]byte("Hello from the media service"))
}

func NewHello(log *zlog.ZwitLogger) *Hello {
	return &Hello{log}
}

type Media struct {
	store             storage.Storage
	config            *config.MediaServiceConfig
	authServiceClient authpb.AuthServiceClient
	Log               *zlog.ZwitLogger
}

type saveFileOutput struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func (m *Media) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	uuid := data.NewUuid()

	m.Log.Infof("Saving a new file: %v", filename)
	err := m.store.Save(fmt.Sprintf("%v/%v", string(uuid), filename), r.Body)
	if err != nil {
		m.Log.Errorf("Failed to save the file to the store: %v", err)
		http.Error(rw, "Failed to save the file.", http.StatusBadRequest)
		return
	}
	m.Log.Infof("Saved the new file to %v", filename)
	rw.Header().Set("Content-Type", "application/json")
	resp := &saveFileOutput{fmt.Sprintf("%v/%v", string(uuid), filename), ""}
	json.NewEncoder(rw).Encode(resp)
}

func NewMedia(s storage.Storage, cfg *config.MediaServiceConfig, authClient authpb.AuthServiceClient, log *zlog.ZwitLogger) *Media {
	return &Media{s, cfg, authClient, log}
}

func setupResponse(rw *http.ResponseWriter, req *http.Request) {
	(*rw).Header().Set("Access-Control-Allow-Origin", "*")
	(*rw).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*rw).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, token")
}

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Created  string `json:"created"`
}

func SetCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		setupResponse(&rw, r)
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(rw, r)
	})
}

func (m *Media) FilesServerLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		m.Log.Info("Files server, getting new file")
		next.ServeHTTP(rw, r)
	})
}

func (m *Media) VerifyTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			m.Log.Errorf("Received an invalid method %v, Only POST is allowed in this endpoint", r.Method)
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		token := r.Header.Get("token")

		resp, err := m.authServiceClient.AuthenticateToken(context.Background(), &authpb.AuthenticateTokenRequest{Token: token})

		if err != nil {
			m.Log.Errorf("Failed to authenticate the user: Could not connect to auth service")
			http.Error(rw, "Failed to authenticate client", http.StatusInternalServerError)
			return
		}

		if !(resp.Auth) {
			m.Log.Errorf("Failed to authenticate the user: invalid token")
			http.Error(rw, "Failed to authenticate client", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUserData{}, &User{Id: resp.GetUser().Id, Username: resp.GetUser().Username})
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}

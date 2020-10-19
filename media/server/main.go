package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dpakach/zwitter/media/api"
	"github.com/dpakach/zwitter/media/storage"
	"github.com/dpakach/zwitter/pkg/auth"
	"github.com/dpakach/zwitter/pkg/config"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.NewMediaSericeConfig("config/config.yaml")
	if err != nil {
		panic(fmt.Errorf("Failed to read config: %s", err))
	}

	err, authEndpoint := cfg.ServiceConfig.GetNodeAddr("Auth")
	if err != nil {
		log.Fatal(err)
	}
	conn, AuthClient := auth.NewAuthClient(authEndpoint)
	defer conn.Close()

	_, err = os.Stat(cfg.LocalStore)
	if err != nil {
		if os.IsNotExist(err) {
			mkdirErr := os.MkdirAll(cfg.LocalStore, os.ModePerm)
			if mkdirErr != nil {
				log.Fatalf(mkdirErr.Error())
			}
		} else {
			log.Fatal(err)
		}
	}

	mediaHandler := api.NewMedia(
		&storage.Local{BasePath: cfg.LocalStore},
		cfg,
		AuthClient,
	)

	helloHandler := api.NewHello()

	sm := mux.NewRouter()

	helloSubrouter := sm.Methods(http.MethodPost).Subrouter()
	helloSubrouter.Handle("/media/ping", helloHandler)

	postHandler := sm.Methods(http.MethodPost).Subrouter()
	postHandler.Handle("/media/{filename:[a-zA-Z]+\\.[a-z]{3}}", mediaHandler)
	postHandler.Use(mediaHandler.VerifyTokenMiddleware)

	getHandler := sm.Methods(http.MethodGet).Subrouter()
	getHandler.Handle(
		"/media/{uuid:[0-9a-f\\-]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/media", http.FileServer(http.Dir(cfg.LocalStore))),
	)

	// create a new server
	s := http.Server{
		Addr:    cfg.Server.RestAddr,       // configure the bind address
		Handler: api.SetCorsMiddleware(sm), // set the default handler
	}

	// start the server
	go func() {
		fmt.Println("Starting server on " + cfg.Server.RestAddr)

		err := s.ListenAndServe()
		if err != nil {
			fmt.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	fmt.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}

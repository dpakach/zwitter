package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dpakach/zwitter/pkg/config"
	"github.com/dpakach/zwitter/web/api"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.NewWebSericeConfig("config/config.yaml")
	if err != nil {
		panic(fmt.Errorf("Failed to read config: %s", err))
	}

	_, err = os.Stat(cfg.AssetsPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf(err.Error())
		} else {
			log.Fatal(err)
		}
	}

	webHandler := api.NewWeb(cfg, "index.html")

	helloHandler := api.NewHello()

	sm := mux.NewRouter().StrictSlash(true)

	helloSubrouter := sm.Methods(http.MethodPost).Subrouter()
	helloSubrouter.Handle("/web/ping", helloHandler)

	getHandler := sm.Methods(http.MethodGet).Subrouter()
	getHandler.PathPrefix("/").Handler(webHandler)

	// create a new server
	s := http.Server{
		Addr:    cfg.Server.RestAddr, // configure the bind address
		Handler: sm,                  // set the default handler
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

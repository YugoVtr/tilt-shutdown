package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	graceful "github.com/yugovtr/tilt-shutdown/http"
	"github.com/yugovtr/tilt-shutdown/mux"
)

func main() {
	const timeout = 30
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT required")
	}
	logger := slog.Default().With("service", "api")

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux.M{
			Duration: timeout,
			Logger:   logger,
		}.Mux(),
		ReadTimeout:  timeout * time.Second,
		WriteTimeout: timeout * time.Second,
	}

	gracefulServer := graceful.NewServer(server, graceful.WithLogger(logger))
	if err := gracefulServer.ListenAndServe(context.Background()); err != nil {
		log.Fatal(err)
	}
}

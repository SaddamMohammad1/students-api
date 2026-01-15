package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SaddamMohammad1/students-api/internal/config"
	"github.com/SaddamMohammad1/students-api/internal/http/handlers/student"
)

// Command to run this project
// go run cmd/students-api/main.go -config config/local.yaml
// here -config is a flag
// config/local.yaml is config file path where present the server related details

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup
	// setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())
	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.Addr))
	fmt.Printf("server started %s", cfg.Addr)

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shurdown successfully")

}

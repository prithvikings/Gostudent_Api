package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prithvikings/Gostudent_Api/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup

	// setup router
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to student API!"))
	})

	// setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("Starting server", slog.String("addr", cfg.Addr))

	// graceful shutdown
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-done
	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server exited properly")
}

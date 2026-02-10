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
	"github.com/prithvikings/Gostudent_Api/internal/http/handlers/student"
	"github.com/prithvikings/Gostudent_Api/internal/storage/sqlite"
)

func main() {
	// ----------------------------
	// Load config
	// ----------------------------
	cfg := config.MustLoad()

	// ----------------------------
	// Database setup
	// ----------------------------
	_, err := sqlite.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	slog.Info("Database initialized successfully")

	// ----------------------------
	// Setup router
	// ----------------------------
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New())

	// ----------------------------
	// Setup server
	// ----------------------------
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Starting server...", slog.String("addr", cfg.Addr))

	// ----------------------------
	// Graceful shutdown channel
	// ----------------------------
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// ----------------------------
	// Start server in goroutine
	// ----------------------------
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// ----------------------------
	// Wait for shutdown signal
	// ----------------------------
	<-done
	slog.Info("Shutting down server...")

	// ----------------------------
	// Shutdown with timeout
	// ----------------------------
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server exited properly")
}

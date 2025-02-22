package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"posts/application"
	"posts/infrastructure"
	"posts/interfaces"
	"posts/internal"
	"posts/utils"
	"syscall"
	"time"

	_ "github.com/lib/pq" // Replace with the appropriate driver for your database
)

var PORT = ":8080"

func main() {
	port := os.Getenv("PORT")
	if port != "" {
		PORT = fmt.Sprintf(":%v", port)
	}
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	slog.Info("init db...")
	db := internal.InitPostgres()
	defer func() {
		slog.Info("Closing PostgreSQL connection...")
		db.Close()
	}()

	slog.Info("init redis...")
	rdb := internal.InitRedis()
	defer func() {
		slog.Info("Closing Redis connection...")
		rdb.Close()
	}()

	cache := infrastructure.NewRedisCache(rdb)

	postRepo := infrastructure.NewPostRepository(db)
	postService := application.NewPostService(postRepo)
	postHandler := interfaces.NewPostHTTPHandler(postService)

	// Create a custom router
	router := utils.NewRouter()

	// Define routes
	router.HandleFunc("GET /{id}", interfaces.LoggerMiddleware(postHandler.GetPost))
	router.HandleFunc("PUT /{id}", interfaces.LoggerMiddleware(postHandler.UpdatePost))
	router.HandleFunc("DELETE /{id}", interfaces.LoggerMiddleware(postHandler.DeletePost))
	router.HandleFunc("POST /", interfaces.LoggerMiddleware(postHandler.CreatePost))
	// this is mocked. Implement soft delete. make visibility = none

	// Start server
	server := &http.Server{Addr: PORT, Handler: router}

	// Start server in a goroutine
	go func() {
		slog.Info(fmt.Sprintf("Starting server on %v", PORT))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // Wait for termination signal
	slog.Info("Shutting down server...")

	// Gracefully shutdown the HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down server", "error", err)
	} else {
		slog.Info("Server shut down successfully")
	}
}

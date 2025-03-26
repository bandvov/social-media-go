package main

import (
	"comments/application"
	"comments/infrastructure"
	"comments/interfaces"
	"comments/internal"
	"comments/utils"
	"context"
	"fmt"
	"log"

	"log/slog"
	"net/http"
	"os"
	"os/signal"
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

	tp := internal.InitTracer("users-service")
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// client := &http.Client{
	// 	Transport: otelhttp.NewTransport(http.DefaultTransport),
	// }
	tracer := tp.Tracer("users-tracer")

	cache := infrastructure.NewRedisCache(rdb)

	// Initialize the CommentFetcher
	commentFetcher := application.NewCommentFetcher(&http.Client{})

	commentRepo := infrastructure.NewPostgresCommentRepository(db, cache)
	commentService := application.NewCommentService(commentRepo, commentFetcher, tracer)
	commentHandler := interfaces.NewCommentHandler(commentService, db, rdb, tracer)

	// Create a custom router
	router := utils.NewRouter()

	router.HandleFunc("POST /", interfaces.LoggerMiddleware(commentHandler.AddComment))
	router.HandleFunc("GET /{id}", interfaces.LoggerMiddleware(commentHandler.GetCommentsByEntityID))
	router.HandleFunc("GET /count", interfaces.LoggerMiddleware(commentHandler.GetCommentsAndRepliesCount))

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

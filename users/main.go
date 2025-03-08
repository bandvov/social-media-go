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
	"users/application"
	"users/infrastructure"
	"users/interfaces"
	internal "users/internal"
	"users/utils"

	_ "github.com/lib/pq" // Replace with the appropriate driver for your database
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

	// Initialize PostgreSQL repository
	userRepo := infrastructure.NewUserRepository(db, cache)

	// Initialize service
	userService := application.NewUserService(userRepo, tracer)

	// Initialize HTTP handler
	userHandler := interfaces.NewUserHTTPHandler(userService, db, rdb, tracer)

	// Create a custom router
	router := utils.NewRouter()

	// seeds.Seed(db, "./migrations/create_users_table.sql")
	// seeds.Seed(db, "./seeds/seed_users.sql")

	// Define routes
	// router.HandleFunc("/api/admin/users", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.IsAdminMiddleware(userHandler.GetAdminProfiles))))

	router.HandleFunc("GET /", interfaces.LoggerMiddleware(userHandler.GetPublicProfiles))
	router.HandleFunc("POST /", interfaces.LoggerMiddleware(userHandler.RegisterUser))
	router.HandleFunc("PUT /{id}", interfaces.LoggerMiddleware(userHandler.UpdateUser))
	router.HandleFunc("GET /{id}/profile", interfaces.LoggerMiddleware(userHandler.GetUserProfile))
	router.HandleFunc("PUT /{id}/role", interfaces.LoggerMiddleware(userHandler.ChangeUserRole))
	router.HandleFunc("POST /login", interfaces.LoggerMiddleware(userHandler.Login))
	router.HandleFunc("POST /by-ids", interfaces.LoggerMiddleware(userHandler.GetUsersByIDs))
	router.HandleFunc("GET /healthz", interfaces.LoggerMiddleware(userHandler.HealthCheckHandler))
	// router.HandleFunc("/seed", seeds.SeedData(db))

	// Start server
	server := &http.Server{Addr: PORT, Handler: otelhttp.NewHandler(router, "users")}

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

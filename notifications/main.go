package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"n/application"
	infrastructure "n/infrastucture"
	"n/interfaces"
	"n/middlewares"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	pgUser := os.Getenv("POSTGRES_USER")
	if pgUser == "" {
		log.Fatal("POSTGRES_USER is not set in the environment")
	}
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	if pgPassword == "" {
		log.Fatal("POSTGRES_PASSWORD is not set in the environment")
	}
	pgDb := os.Getenv("POSTGRES_DB")
	if pgDb == "" {
		log.Fatal("POSTGRES_DB is not set in the environment")
	}
	pgPort := os.Getenv("POSTGRES_PORT")
	if pgPort == "" {
		log.Fatal("POSTGRES_PORT is not set in the environment")
	}
	dbConnectionString := fmt.Sprintf("postgresql://%v:%v@localhost:%v/%v?sslmode=disable", pgUser, pgPassword, pgPort, pgDb)

	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Redis connection
	redis := infrastructure.NewRedisEventListener("localhost:6379")

	// Dependency Injection
	repo := infrastructure.NewPostgresNotificationRepository(db)
	service := application.NewNotificationService(repo, redis)
	handler := interfaces.NewNotificationHandler(service)

	// HTTP Router
	r := http.NewServeMux()

	r.HandleFunc("/", handler.GetNotifications)
	r.HandleFunc("/send", handler.SendNotification)
	r.HandleFunc("/listen", handler.ListenNotifications)
	r.HandleFunc("/mark_as_read", handler.MarkAsRead)
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// Check the database connection
		if err := db.Ping(); err != nil {
			http.Error(w, "Database unreachable", http.StatusInternalServerError)
			return
		}

		// Check the Redis connection
		if err := redis.Ping(context.Background()).Err(); err != nil {
			http.Error(w, "Redis unreachable", http.StatusInternalServerError)
			return
		}

		// Everything is OK, return a 200 status
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", middlewares.CorsMiddleware(r)))
}

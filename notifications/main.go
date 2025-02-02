package main

import (
	"database/sql"
	"fmt"
	"log"
	"n/application"
	infrastructure "n/infrastucture"
	"n/interfaces"
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

	r.HandleFunc("/send", handler.SendNotification)
	r.HandleFunc("/listen", handler.ListenNotifications)

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

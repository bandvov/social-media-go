package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bandvov/social-media-go/application"
	"github.com/bandvov/social-media-go/infrastructure"
	"github.com/bandvov/social-media-go/interfaces"
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
	dbConnectionString := fmt.Sprintf("postgresql://%v:%v@localhost:%v/%v?sslmode=disable", pgUser, pgPassword, pgDb, pgPort)

	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize PostgreSQL repository
	repo := infrastructure.NewPostgresUserRepository(db)

	// Initialize service
	userService := application.NewUserService(repo)

	// Initialize HTTP handler
	handler := interfaces.NewHTTPHandler(userService)

	// Set up routes
	http.HandleFunc("/register", handler.RegisterUser)
	http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			handler.UpdateUser(w, r)
		} else if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/role") {
			handler.ChangeUserRole(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	// Start server
	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

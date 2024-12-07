package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bandvov/social-media-go/application"
	"github.com/bandvov/social-media-go/infrastructure"
	"github.com/bandvov/social-media-go/interfaces"
	"github.com/bandvov/social-media-go/seeds"
	"github.com/bandvov/social-media-go/utils"
	_ "github.com/lib/pq" // Replace with the appropriate driver for your database
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

	// Initialize PostgreSQL repository
	repo := infrastructure.NewPostgresUserRepository(db)

	// Initialize service
	userService := application.NewUserService(repo)

	// Initialize HTTP handler
	handler := interfaces.NewHTTPHandler(userService)

	// Create a custom router
	router := utils.NewRouter()

	seeds.Seed(db, "./migrations/create_users_table.sql")
	seeds.Seed(db, "./migrations/create_posts_table.sql")

	seeds.Seed(db, "./seeds/seed_users.sql")
	seeds.Seed(db, "./seeds/seed_posts.sql")

	// Define routes
	router.Handle("POST", "/register", interfaces.LoggerMiddleware(handler.RegisterUser))
	router.Handle("PUT", "/user/", interfaces.LoggerMiddleware(handler.AuthMiddleware(handler.UpdateUser)))
	router.Handle("POST", "/user/role", interfaces.LoggerMiddleware(handler.AuthMiddleware(handler.ChangeUserRole)))
	router.Handle("POST", "/login", interfaces.LoggerMiddleware(handler.Login))

	// Start server
	log.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

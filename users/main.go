package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"users/application"
	"users/infrastructure"
	"users/interfaces"
	"users/utils"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq" // Replace with the appropriate driver for your database
)

var PORT = ":8080"

func main() {
	port := os.Getenv("PORT")
	if port != "" {
		PORT = fmt.Sprintf(":%v", port)
	}

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

	// Redis setup
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cache := infrastructure.NewRedisCache(redisClient)

	// Initialize PostgreSQL repository
	userRepo := infrastructure.NewUserRepository(db, cache)

	// Initialize service
	userService := application.NewUserService(userRepo)

	// Initialize HTTP handler
	userHandler := interfaces.NewUserHTTPHandler(userService)

	// Create a custom router
	router := utils.NewRouter()

	// seeds.Seed(db, "./migrations/create_users_table.sql")
	// seeds.Seed(db, "./seeds/seed_users.sql")
	

	// Define routes
	// router.HandleFunc("/api/admin/users", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.IsAdminMiddleware(userHandler.GetAdminProfiles))))

	router.HandleFunc("GET /api/users", interfaces.LoggerMiddleware(userHandler.GetPublicProfiles))
	router.HandleFunc("POST /api/users", interfaces.LoggerMiddleware(userHandler.RegisterUser))
	router.HandleFunc("PUT /api/users/{id}", interfaces.LoggerMiddleware(userHandler.UpdateUser))
	router.HandleFunc("GET /api/users/{id}/profile", interfaces.LoggerMiddleware(userHandler.GetUserProfile))
	router.HandleFunc("PUT /api/users/{id}/role", interfaces.LoggerMiddleware(userHandler.ChangeUserRole))
	router.HandleFunc("POST /api/users/login", interfaces.LoggerMiddleware(userHandler.Login))
	// router.HandleFunc("/seed", seeds.SeedData(db))

	// Start server
	log.Printf("Server is running on %v", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}

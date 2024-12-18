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
	"github.com/bandvov/social-media-go/utils"
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
	// Initialize PostgreSQL repository
	userRepo := infrastructure.NewUserRepository(db)

	// Initialize service
	userService := application.NewUserService(userRepo)

	// Initialize HTTP handler
	userHandler := interfaces.NewUserHTTPHandler(userService)

	postRepo := infrastructure.NewPostRepository(db)
	postService := application.NewPostService(postRepo)
	postHandler := interfaces.NewPostHTTPHandler(postService)

	// Create a custom router
	router := utils.NewRouter()

	// seeds.Seed(db, "./migrations/create_users_table.sql")
	// seeds.Seed(db, "./migrations/create_posts_table.sql")
	// seeds.Seed(db, "./migrations/media_urls_create_table.sql")
	// seeds.Seed(db, "./migrations/create_reactions_table.sql")
	// seeds.Seed(db, "./migrations/create_reaction_types.table.sql")

	// seeds.Seed(db, "./seeds/seed_users.sql")
	// seeds.Seed(db, "./seeds/seed_posts.sql")
	// seeds.Seed(db, "./seeds/seed_media_urls.sql")
	// seeds.Seed(db, "./seeds/seed_reaction_types.sql")
	// seeds.Seed(db, "./seeds/seed_reactions.sql")

	// Define routes
	router.HandleFunc("GET /user", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.GetUserProfile)))
	router.HandleFunc("POST /user", interfaces.LoggerMiddleware(userHandler.RegisterUser))
	router.HandleFunc("DELETE /user/", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.UpdateUser)))
	router.HandleFunc("POST /user/login", interfaces.LoggerMiddleware(userHandler.Login))
	router.HandleFunc("PUT /user/role", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.ChangeUserRole)))
	router.HandleFunc("GET /user/all", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.GetAllUsers)))

	router.HandleFunc("GET /post", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.Get)))
	router.HandleFunc("POST /post", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.Create)))
	router.HandleFunc("PUT /post", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.Update)))
	router.HandleFunc("DELETE /post", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.Delete)))

	// Start server
	log.Printf("Server is running on %v", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}

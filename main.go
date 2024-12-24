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

	followerRepo := infrastructure.NewFollowerRepository(db)
	Followerservice := application.NewFollowerService(followerRepo)
	followerHandler := interfaces.NewFollowerHandler(Followerservice)

	tagRepo := infrastructure.NewPostgresTagRepository(db)
	tagService := application.NewTagService(tagRepo)
	tagHandler := interfaces.NewTagHandler(tagService)

	// Create a custom router
	router := utils.NewRouter()

	// seeds.Seed(db, "./migrations/create_users_table.sql")
	// seeds.Seed(db, "./migrations/create_posts_table.sql")
	// seeds.Seed(db, "./migrations/media_urls_create_table.sql")
	// seeds.Seed(db, "./migrations/create_reactions_table.sql")
	// seeds.Seed(db, "./migrations/create_reaction_types.table.sql")
	// seeds.Seed(db, "./migrations/create_followers_table.sql")

	// seeds.Seed(db, "./seeds/seed_users.sql")
	// seeds.Seed(db, "./seeds/seed_posts.sql")
	// seeds.Seed(db, "./seeds/seed_media_urls.sql")
	// seeds.Seed(db, "./seeds/seed_reaction_types.sql")
	// seeds.Seed(db, "./seeds/seed_reactions.sql")
	// seeds.Seed(db, "./seeds/seed_followers.sql")

	// Define routes
	router.HandleFunc("GET /users/{id}/profile", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.GetUserProfile)))
	router.HandleFunc("POST /users", interfaces.LoggerMiddleware(userHandler.RegisterUser))
	router.HandleFunc("PUT /users/{id}", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.UpdateUser)))
	router.HandleFunc("POST /users/login", interfaces.LoggerMiddleware(userHandler.Login))
	router.HandleFunc("PUT /users/{id}/role", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.ChangeUserRole)))
	router.HandleFunc("GET /users", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.GetAllUsers)))

	router.HandleFunc("GET /users/{id}/posts", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.GetPostsByUser)))

	router.HandleFunc("GET /posts/{id}", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.GetPost)))
	router.HandleFunc("POST /posts", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.CreatePost)))
	router.HandleFunc("PUT /posts/{id}", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.UpdatePost)))
	// this is mocked. Implement soft delete. make visibility = none
	router.HandleFunc("DELETE /posts/{id}", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.DeletePost)))

	http.HandleFunc("POST /followers", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(followerHandler.AddFollower)))
	http.HandleFunc("DELETE /followers/{id}", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(followerHandler.RemoveFollower)))
	http.HandleFunc("GET /followers", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(followerHandler.GetFollowers)))

	http.HandleFunc("GET /tags", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(tagHandler.GetTags)))
	http.HandleFunc("POST /tags", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(tagHandler.CreateTag)))
	// This in mocked
	http.HandleFunc("DELETE /tags{id}", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(tagHandler.DeleteTag)))
	// Start server
	log.Printf("Server is running on %v", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}

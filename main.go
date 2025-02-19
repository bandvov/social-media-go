package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bandvov/social-media-go/application"
	"github.com/bandvov/social-media-go/infrastructure"
	"github.com/bandvov/social-media-go/interfaces"
	"github.com/bandvov/social-media-go/utils"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq" // Replace with the appropriate driver for your database
)

var PORT = ":443"

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

	commentRepo := infrastructure.NewPostgresCommentRepository(db)
	commentService := application.NewCommentService(commentRepo)
	commentHandler := interfaces.NewCommentHandler(commentService)

	reactionRepo := infrastructure.NewReactionRepository(db)
	reactionService := application.NewReactionService(reactionRepo)
	reactionHandler := interfaces.NewReactionHandler(reactionService)

	postRepo := infrastructure.NewPostRepository(db)
	postService := application.NewPostService(postRepo)
	postHandler := interfaces.NewPostHTTPHandler(postService, commentService, userService, reactionService)

	followerRepo := infrastructure.NewFollowerRepository(db)
	Followerservice := application.NewFollowerService(followerRepo)
	followerHandler := interfaces.NewFollowerHandler(Followerservice)

	tagRepo := infrastructure.NewTagRepository(db)
	tagService := application.NewTagService(tagRepo)
	tagHandler := interfaces.NewTagHandler(tagService)

	// Create a custom router
	router := utils.NewRouter()

	// seeds.Seed(db, "./migrations/create_users_table.sql")
	// seeds.Seed(db, "./migrations/create_posts_table.sql")
	// seeds.Seed(db, "./migrations/media_urls_create_table.sql")
	// seeds.Seed(db, "./migrations/create_reaction_types.table.sql")
	// seeds.Seed(db, "./migrations/create_reactions_table.sql")
	// seeds.Seed(db, "./migrations/create_followers_table.sql")
	// seeds.Seed(db, "./migrations/create_tags_table.sql")
	// seeds.Seed(db, "./migrations/create_comments_table.sql")

	// seeds.Seed(db, "./seeds/seed_users.sql")
	// seeds.Seed(db, "./seeds/seed_posts.sql")
	// seeds.Seed(db, "./seeds/seed_media_urls.sql")
	// seeds.Seed(db, "./seeds/seed_reaction_types.sql")
	// seeds.Seed(db, "./seeds/seed_reactions.sql")
	// seeds.Seed(db, "./seeds/seed_followers.sql")
	// seeds.Seed(db, "./seeds/seed_tags.sql")
	// seeds.Seed(db, "./seeds/seed_comments.sql")

	// Define routes
	router.HandleFunc("/api/admin/users", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.IsAdminMiddleware(userHandler.GetAdminProfiles))))

	router.HandleFunc("GET /api/users", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.GetPublicProfiles)))
	router.HandleFunc("GET /api/users/{id}/profile", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.GetUserProfile)))
	router.HandleFunc("POST /api/users", interfaces.LoggerMiddleware(userHandler.RegisterUser))
	router.HandleFunc("PUT /api/users/{id}", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.UpdateUser)))
	router.HandleFunc("POST /api/users/login", interfaces.LoggerMiddleware(userHandler.Login))
	router.HandleFunc("PUT /api/users/{id}/role", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(userHandler.ChangeUserRole)))

	router.HandleFunc("GET /api/users/{id}/posts", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.GetPostsByUser)))
	router.HandleFunc("GET /api/users/{id}/followers", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(followerHandler.GetFollowers)))
	router.HandleFunc("GET /api/users/{id}/followees", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(followerHandler.GetFollowees)))

	router.HandleFunc("GET /api/posts/{id}", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.GetPost)))
	router.HandleFunc("POST /api/posts", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.CreatePost)))
	router.HandleFunc("PUT /api/posts/{id}", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.UpdatePost)))
	// this is mocked. Implement soft delete. make visibility = none
	router.HandleFunc("DELETE /api/posts/{id}", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(postHandler.DeletePost)))

	router.HandleFunc("POST /api/followers", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(followerHandler.AddFollower)))
	router.HandleFunc("DELETE /api/followers/{id}", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(followerHandler.RemoveFollower)))

	router.HandleFunc("GET /tags", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(tagHandler.GetTags)))
	router.HandleFunc("POST /tags", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(tagHandler.CreateTag)))
	router.HandleFunc("DELETE /tags/{id}", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(tagHandler.DeleteTag)))

	router.HandleFunc("POST /api/comments", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(commentHandler.AddComment)))
	router.HandleFunc("GET /api/comments/{id}", interfaces.LoggerMiddleware(userHandler.AuthMiddleware(commentHandler.GetCommentsByEntityID)))

	router.HandleFunc("GET /api/reaction", reactionHandler.AddOrUpdateReaction)
	router.HandleFunc("DELETE /api/reaction", reactionHandler.RemoveReaction)

	// router.HandleFunc("/seed", seeds.SeedData(db))

	// Configure TLS
	tlsConfig := &tls.Config{}

	server := &http.Server{
		Addr:      PORT,
		Handler:   interfaces.CorsMiddleware(router),
		TLSConfig: tlsConfig,
	}

	// Start server
	log.Printf("Server is running on %v", PORT)
	log.Fatal(server.ListenAndServeTLS("./certs/cert.pem", "./certs/key.pem"))
}

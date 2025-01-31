package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bandvov/social-media-go/auth/application"
	"github.com/bandvov/social-media-go/auth/infrastructure"
	"github.com/bandvov/social-media-go/auth/interfaces"
)

var PORT = ":8081"

func main() {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatal("No JWT secret key")
	}
	port := os.Getenv("PORT")
	if port == "" {
		PORT = "8081"
	}
	authService := infrastructure.NewJWTAuthService(secretKey)
	authApp := application.NewAuthApplication(authService)
	authHandler := interfaces.NewAuthHandler(authApp)

	http.HandleFunc("/verify", authHandler.VerifyTokenHandler)

	log.Printf("Server running on port%v ...\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

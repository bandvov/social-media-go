package main

import (
	"auth-service/application"
	"auth-service/infrastructure"
	"auth-service/interfaces"
	"log"
	"net/http"
	"os"
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

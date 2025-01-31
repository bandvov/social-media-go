package main

import (
	"activity-service/application"
	"activity-service/infrastructure"
	"activity-service/interfaces"
	"activity-service/utils"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // Replace with the appropriate driver for your database
)

var PORT = ":80"

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

	activityRepo := infrastructure.NewPostgresActivityRepository(db)
	activityService := application.NewActivityService(activityRepo)
	activityHabdler := interfaces.NewActivityHandler(activityService)

	router := utils.NewRouter()

	router.HandleFunc("POST /activities", activityHabdler.AddActivity)
	router.HandleFunc("GET /activities", activityHabdler.GetActivities)
	// Start server
	log.Printf("Server is running on %v", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))

}

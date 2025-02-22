package internal

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"time"
	"users/internal/config"
)

// Initialize PostgreSQL connection
func InitPostgres() *sql.DB {
	cfg := config.LoadConfig()

	connStr := fmt.Sprintf("postgresql://%v:%v@localhost:%v/%v?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresPort, cfg.PostgresDB)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Warn("PostgreSQL health check failed", "error", err)
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return db
}

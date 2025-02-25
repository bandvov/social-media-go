package interfaces

import (
	"log"
	"net/http"
)

// Define keys for context
type contextKey string

const (
	userIDKey  contextKey = "userID"
	isAdminKey contextKey = "isAdmin"
)

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

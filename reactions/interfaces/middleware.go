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

func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assuming admin status is part of context
		isAdmin := r.Context().Value(isAdminKey).(bool)
		if !isAdmin {
			http.Error(w, "forbidden: admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

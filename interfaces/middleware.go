package interfaces

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bandvov/social-media-go/utils"
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
		isAdmin := r.Context().Value("isAdmin").(bool)
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

// Middleware to extract userID from cookie and add to context
func (h *UserHTTPHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the cookie
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse userID from cookie
		var token string
		_, err = fmt.Sscanf(cookie.Value, "%s", &token)
		if err != nil {
			http.Error(w, "Invalid access token", http.StatusBadRequest)
			return
		}

		userID, err := utils.ValidateJWT(token)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Retrieve user from the database
		user, err := h.UserService.GetUserByID(userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		isAdmin := user.Role == "admin"

		// Add userID and isAdmin to context
		ctx := context.WithValue(r.Context(), userIDKey, user.ID)
		ctx = context.WithValue(ctx, isAdminKey, isAdmin)
		// Call the next handler with updated context
		next(w, r.WithContext(ctx))
	}
}

// Middleware to extract userID from cookie and add to context
func (h *UserHTTPHandler) IsAdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve userID and isAdmin from context
		isAdmin := r.Context().Value(isAdminKey).(bool)
		if !isAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		// Call the next handler with updated context
		next(w, r)
	}
}

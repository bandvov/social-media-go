package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/bandvov/social-media-go/auth/application"
)

type AuthHandler struct {
	authApp *application.AuthApplication
}

func NewAuthHandler(authApp *application.AuthApplication) *AuthHandler {
	return &AuthHandler{authApp: authApp}
}

func (h *AuthHandler) VerifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	userID, err := h.authApp.Authenticate(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := map[string]string{"user_id": userID}
	json.NewEncoder(w).Encode(response)
}

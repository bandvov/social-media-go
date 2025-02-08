package interfaces

import (
	"auth-service/application"
	"encoding/json"
	"net/http"
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

	userID, err := h.authApp.ValidateServiceToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	response := map[string]string{"user_id": userID}
	json.NewEncoder(w).Encode(response)
}

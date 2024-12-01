package interfaces

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/bandvov/social-media-go/application"
)

type HTTPHandler struct {
	UserService *application.UserService
}

func NewHTTPHandler(userService *application.UserService) *HTTPHandler {
	return &HTTPHandler{UserService: userService}
}

func parseUserIDFromPath(path string) (int64, error) {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) < 2 {
		return 0, errors.New("invalid URL")
	}
	return strconv.ParseInt(segments[1], 10, 64)
}

func (h *HTTPHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := ValidateEmail(req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ValidatePassword(req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.UserService.RegisterUser(req.Username, req.Password, req.Email, req.FirstName, req.LastName)
	if err != nil {
		http.Error(w, "error registering user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user registered successfully"})
}

func (h *HTTPHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := parseUserIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Bio        string `json:"bio"`
		ProfilePic string `json:"profile_pic"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email != "" {
		if err := ValidateEmail(req.Email); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if req.Password != "" {
		if err := ValidatePassword(req.Password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	err = h.UserService.UpdateUserData(userID, req.Email, req.Password, req.FirstName, req.LastName, req.Bio, req.ProfilePic)
	if err != nil {
		http.Error(w, "error updating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "user updated successfully"})
}

func (h *HTTPHandler) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, err := parseUserIDFromPath(strings.TrimSuffix(r.URL.Path, "/role"))
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		NewRole string `json:"new_role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := ValidateRole(req.NewRole); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isAdmin := r.Context().Value("is_admin").(bool)

	err = h.UserService.ChangeUserRole(userID, req.NewRole, isAdmin)
	if err != nil {
		http.Error(w, "error changing user role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "user role changed successfully"})
}

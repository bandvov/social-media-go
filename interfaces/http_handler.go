package interfaces

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bandvov/social-media-go/application"
	"github.com/bandvov/social-media-go/domain"
	"github.com/bandvov/social-media-go/utils"
)

type HTTPHandler struct {
	UserService application.UserServiceInterface
}

func NewHTTPHandler(userService application.UserServiceInterface) *HTTPHandler {
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
	var newUser domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := ValidateEmail(newUser.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ValidatePassword(newUser.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.UserService.RegisterUser(newUser)
	if err != nil {
		http.Error(w, "error registering user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user registered successfully"})
}

func (h *HTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse and validate the request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if request.Email == "" || request.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, err := h.UserService.Authenticate(request.Email, request.Password)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(int(user.ID), user.Username)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set token in cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	})

	// Respond with user data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id":   user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"role":      user.Role,
	})
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
	fmt.Println("here")
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

func (h *HTTPHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenString := authHeader[7:]
	userID, err := utils.ValidateJWT(tokenString)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Fetch user profile from service
	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with user profile data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

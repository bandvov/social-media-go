package interfaces

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"users/application"
	"users/domain"
	"users/utils"

	"github.com/lib/pq"
)

// Define keys for context
type contextKey string

const (
	userIDKey  contextKey = "userID"
	isAdminKey contextKey = "isAdmin"
)

type UserHTTPHandler struct {
	UserService application.UserServiceInterface
	db          *sql.DB
}

func NewUserHTTPHandler(userService application.UserServiceInterface, db *sql.DB) *UserHTTPHandler {
	return &UserHTTPHandler{UserService: userService, db: db}
}

func (h *UserHTTPHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser struct {
		Data domain.CreateUserRequest `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, `{"message": "invalid request body"}`, http.StatusBadRequest)
		return
	}
	if err := ValidateEmail(newUser.Data.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ValidatePassword(newUser.Data.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.UserService.RegisterUser(newUser.Data)
	if err != nil {
		fmt.Println("err: ", err)
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			http.Error(w, "error registering user: user already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, "error registering user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user registered successfully"})
}

func (h *UserHTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Data domain.CreateUserRequest `json:"data"`
	}

	// Parse and validate the request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "{\"message\": \"invalid request body\"}", http.StatusBadRequest)
		return
	}

	if request.Data.Email == "" || request.Data.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, err := h.UserService.Authenticate(request.Data.Email, request.Data.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	// Set token in cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		Secure:   true,
	})

	// Respond with user data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": user.ID})
}

func (h *UserHTTPHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "{\"message\": \"invalid user ID\"}", http.StatusBadRequest)
		return
	}

	req := &domain.User{}
	req.ID = userID

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, "{\"message\": \"invalid request body\"}", http.StatusBadRequest)
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

	err = h.UserService.UpdateUserData(req)
	if err != nil {
		http.Error(w, "error updating user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "user updated successfully"})
}

func (h *UserHTTPHandler) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", struct{ message string }{message: "invalid user ID"}), http.StatusBadRequest)
		return
	}

	var req struct {
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := ValidateRole(req.Role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isAdmin := r.Context().Value(isAdminKey).(bool)

	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err = h.UserService.ChangeUserRole(userID, req.Role, isAdmin)
	if err != nil {
		http.Error(w, "error changing user role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "user role changed successfully"})
}

func (h *UserHTTPHandler) GetPublicProfiles(w http.ResponseWriter, r *http.Request) {
	limit, offset := utils.ParsePagination(r)
	users, err := h.UserService.GetPublicProfiles(limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch public profiles", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHTTPHandler) GetAdminProfiles(w http.ResponseWriter, r *http.Request) {
	isAdmin := r.Context().Value(isAdminKey).(bool)

	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	limit, offset := utils.ParsePagination(r)
	users, err := h.UserService.GetAdminProfiles(limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch admin profiles", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHTTPHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "{\"message\": \"invalid user ID\"}", http.StatusBadRequest)
		return
	}

	// Ensure user lookup happens after authorization checks
	user, err := h.UserService.GetUserProfileInfo(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user.Password = ""
	// Respond with user profile data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHTTPHandler) IsAdmin(ctx context.Context) bool {
	return ctx.Value(isAdminKey).(bool)
}

func (h *UserHTTPHandler) Verify(w http.ResponseWriter, r *http.Request) {
	cookieName := "access_token"
	// Extract the cookie
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Println("here========================")
	// Parse userID from cookie
	var token string
	_, err = fmt.Sscanf(cookie.Value, "%s", &token)
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusBadRequest)
		return
	}

	fmt.Println("here1========================")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Path:     "/",
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Unix(0, 0),
		})
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Println("here2========================")
	// Retrieve user from the database
	user, err := h.UserService.GetUserByID(claims.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	fmt.Println("here4========================")

	isAdmin := user.Role == "admin"

	response := map[string]interface{}{
		string(userIDKey):  user.ID,
		string(isAdminKey): isAdmin,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHTTPHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		slog.Warn("PostgreSQL health check failed", "error", err)
		http.Error(w, "Unhealthy", http.StatusServiceUnavailable)
		return
	}

	slog.Info("Health check passed")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Ok")
}

package interfaces

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bandvov/social-media-go/application"
	"github.com/bandvov/social-media-go/domain"
	"github.com/bandvov/social-media-go/utils"
	"github.com/lib/pq"
)

type UserHTTPHandler struct {
	UserService application.UserServiceInterface
}

func NewUserHTTPHandler(userService application.UserServiceInterface) *UserHTTPHandler {
	return &UserHTTPHandler{UserService: userService}
}

func parseUserIDFromPath(path string) (int, error) {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) < 2 {
		return 0, errors.New("invalid URL")
	}
	userId, err := strconv.ParseInt(segments[1], 10, 64)
	if err != nil {
		return 0, err
	}
	return int(userId), nil
}

func (h *UserHTTPHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "{\"message\": \"invalid request body\"}", http.StatusBadRequest)
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
	var request domain.CreateUserRequest

	// Parse and validate the request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "{\"message\": \"invalid request body\"}", http.StatusBadRequest)
		return
	}
	if request.Email == "" || request.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, err := h.UserService.Authenticate(request.Email, request.Password)
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
		HttpOnly: true,
		Path:     "/",
	})

	// Respond with user data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
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

	var req domain.User
	req.ID = userID

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

func (h *UserHTTPHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(userIDKey).(interface{}).(int)
	if !ok || userId == 0 {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	isAdmin := h.IsAdmin(r.Context())

	id := r.PathValue("id")
	userIDFromUrl, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
	if !isAdmin || userId != userIDFromUrl {
		fmt.Println("here1")
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Fetch user profile from service
	user, err := h.UserService.GetUserByID(userIDFromUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user.Password = ""
	fmt.Println("here2")
	// Respond with user profile data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHTTPHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Parse `limit` and `offset` with default values
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // Default limit
	}

	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0 // Default offset
	}

	// Parse `sort` with default value
	sort := query.Get("sort")
	if sort != "asc" && sort != "desc" {
		sort = "desc" // Default sort
	}
	search := query.Get("search")
	orderBy := query.Get("order_by")

	users, err := h.UserService.GetAllUsers(limit, offset, sort, orderBy, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHTTPHandler) IsAdmin(ctx context.Context) bool {
	return ctx.Value(isAdminKey).(bool)
}

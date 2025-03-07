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
	"users/internal"
	"users/utils"

	"github.com/lib/pq"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.22.0"
	"go.opentelemetry.io/otel/trace"
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
	rdb         internal.RedisClient
	tracer      trace.Tracer
}

func NewUserHTTPHandler(userService application.UserServiceInterface, db *sql.DB, rdb internal.RedisClient, tracer trace.Tracer) *UserHTTPHandler {
	return &UserHTTPHandler{UserService: userService, db: db, rdb: rdb, tracer: tracer}
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

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err := h.UserService.RegisterUser(ctx, newUser.Data)
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

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Authenticate user
	user, err := h.UserService.Authenticate(ctx, request.Data.Email, request.Data.Password)
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

	var req struct {
		Data domain.User `json:"data"`
	}
	req.Data.ID = userID

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "{\"message\": \"invalid request body\"}", http.StatusBadRequest)
		return
	}

	if req.Data.Email != "" {
		if err := ValidateEmail(req.Data.Email); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if req.Data.Password != "" {
		if err := ValidatePassword(req.Data.Password); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	// we don't want to update the role
	// rele is updated by the admin in a different endpoint
	req.Data.Role = ""

	err = h.UserService.UpdateUserData(ctx, &req.Data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error updating user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "user updated successfully"})
}

func (h *UserHTTPHandler) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	headerValue := r.Header.Get("is_admin")

	// Convert string to boolean
	isAdmin, err := strconv.ParseBool(headerValue)
	if err != nil {
		http.Error(w, "Invalid boolean header value", http.StatusBadRequest)
		return
	}

	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id := r.PathValue("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", struct{ message string }{message: "invalid user ID"}), http.StatusBadRequest)
		return
	}

	var req struct {
		Data struct {
			Role string `json:"role"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := ValidateRole(req.Data.Role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err = h.UserService.ChangeUserRole(ctx, userID, req.Data.Role)
	if err != nil {
		http.Error(w, "error changing user role: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "user role changed successfully"})
}

func (h *UserHTTPHandler) GetPublicProfiles(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	p := utils.ParsePagination(r)

	users, err := h.UserService.GetPublicProfiles(ctx, p)
	if err != nil {
		http.Error(w, "Failed to fetch public profiles", http.StatusInternalServerError)
		return
	}

	// use users count to determine if there are more profiles
	response := map[string]interface{}{
		"data":    users,
		"hasMore": 20 > p.Offset+p.Limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHTTPHandler) GetAdminProfiles(w http.ResponseWriter, r *http.Request) {
	isAdmin := r.Context().Value(isAdminKey).(bool)

	if !isAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	p := utils.ParsePagination(r)
	users, err := h.UserService.GetAdminProfiles(ctx, p)
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

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Ensure user lookup happens after authorization checks
	user, err := h.UserService.GetUserProfileInfo(ctx, userID)
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

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	fmt.Println("here2========================")
	// Retrieve user from the database
	user, err := h.UserService.GetUserByID(ctx, claims.UserID)
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

// Handler for getting users by IDs
func (h *UserHTTPHandler) GetUsersByIDs(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	// Parent span
	ctx, parentSpan := h.tracer.Start(ctx, "handler.GetUsersByIDs")
	defer parentSpan.End()
	// Add an annotation when the request is received
	parentSpan.AddEvent("Request received", trace.WithAttributes(
		semconv.HTTPMethodKey.String(r.Method),
		semconv.HTTPURLKey.String(r.URL.String()),
		semconv.HTTPStatusCodeKey.Int(http.StatusOK),
	))

	// Decode the JSON request body
	var request struct {
		Data []int `json:"data"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		parentSpan.AddEvent("unmarshal-error", trace.WithAttributes(
			attribute.Bool("success", false),
			attribute.String("error", err.Error()),
		))
		http.Error(w, fmt.Sprintf("Error parsing request body: %v", err), http.StatusBadRequest)
		return
	}
	parentSpan.AddEvent("unmarshal-success", trace.WithAttributes(
		attribute.Bool("success", true),
	))

	// Call GetUsersByIDs function
	users, err := h.UserService.GetUsersByIDs(ctx, request.Data)
	if err != nil {
		parentSpan.AddEvent("Error fetching users from service", trace.WithAttributes(
			attribute.String("error:", err.Error()),
		))
		http.Error(w, fmt.Sprintf("Error fetching users: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the fetched user data in JSON format
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string][]domain.User{"data": users})
	if err != nil {
		parentSpan.AddEvent("Error marshalling users", trace.WithAttributes(
			attribute.String("error:", err.Error()),
		))
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
	parentSpan.AddEvent("success response", trace.WithAttributes(
		attribute.Bool("success", true),
	))
}

func (h *UserHTTPHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		slog.Warn("PostgreSQL health check failed", "error", err)
		http.Error(w, "Unhealthy", http.StatusServiceUnavailable)
		return
	}
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		slog.Warn("Redis health check failed", "error", err)
		http.Error(w, "Unhealthy", http.StatusServiceUnavailable)
		return
	}

	slog.Info("Health check passed")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Ok")
}

package interfaces

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"followers/application"
	"followers/internal"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type FollowHandler struct {
	service application.FollowServiceInterface
	db      *sql.DB
	rdb     internal.RedisClient
}

func NewFollowHandler(service application.FollowServiceInterface, db *sql.DB, rdb internal.RedisClient) *FollowHandler {
	return &FollowHandler{service: service, db: db, rdb: rdb}
}

func (h *FollowHandler) AddFollower(w http.ResponseWriter, r *http.Request) {
	// Parse the URL parameters to get the follower and followee IDs
	userID, ok := r.Context().Value(userIDKey).(interface{}).(int)
	if !ok || userID == 0 {
		http.Error(w, "unauthenticated", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	followeeID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid follower ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Call the service to add the follower
	err = h.service.AddFollower(ctx, userID, followeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send a response back
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Follower added successfully")
}

func (h *FollowHandler) RemoveFollower(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	targetUserId := query.Get("target_id")
	targetUserIDFromUrl, err := strconv.Atoi(targetUserId)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	followeeID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid followee ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Call the service to remove the follower
	err = h.service.RemoveFollower(ctx, targetUserIDFromUrl, followeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send a response back
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Follower removed successfully")
}

func (h *FollowHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userIDFromUrl, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
	query := r.URL.Query()
	targetUserId := query.Get("target_id")
	targetUserIDFromUrl, err := strconv.Atoi(targetUserId)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

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

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Call the service to get followers
	followers, err := h.service.GetFollowers(ctx, userIDFromUrl, targetUserIDFromUrl, limit, offset, sort, orderBy, search)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}
func (h *FollowHandler) GetFollowees(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userIDFromUrl, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	targetUserId := query.Get("target_id")
	targetUserIDFromUrl, err := strconv.Atoi(targetUserId)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
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

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Call the service to get followers
	followers, err := h.service.GetFollowees(ctx, userIDFromUrl, targetUserIDFromUrl, limit, offset, sort, orderBy, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func (h *FollowHandler) GetFollowerStats(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userIDFromUrl, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	followersCount, followeesCount, err := h.service.GetFollowerStats(ctx, userIDFromUrl)
	if err != nil {
		http.Error(w, "Error fetching follower stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"followers_count": followersCount,
		"followees_count": followeesCount,
	})
}

func (h *FollowHandler) CheckFollowStatus(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	targetUserID, err2 := strconv.ParseInt(r.URL.Query().Get("target_user_id"), 10, 64)
	if err != nil || err2 != nil {
		http.Error(w, "Invalid user_id or target_user_id", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	result, err := h.service.CheckFollowStatus(ctx, userID, targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (h *FollowHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
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

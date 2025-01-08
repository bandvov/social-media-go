package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bandvov/social-media-go/application"
)

type FollowerHandler struct {
	service application.FollowerServiceInterface
}

func NewFollowerHandler(service application.FollowerServiceInterface) *FollowerHandler {
	return &FollowerHandler{service: service}
}

func (h *FollowerHandler) AddFollower(w http.ResponseWriter, r *http.Request) {
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

	// Call the service to add the follower
	err = h.service.AddFollower(userID, followeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send a response back
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Follower added successfully")
}

func (h *FollowerHandler) RemoveFollower(w http.ResponseWriter, r *http.Request) {
	// Parse the URL parameters to get the follower and followee IDs
	userID, ok := r.Context().Value(userIDKey).(interface{}).(int)
	if !ok || userID == 0 {
		http.Error(w, "unauthenticated", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	followeeID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid followee ID", http.StatusBadRequest)
		return
	}

	// Call the service to remove the follower
	err = h.service.RemoveFollower(userID, followeeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send a response back
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Follower removed successfully")
}

func (h *FollowerHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(userIDKey).(interface{}).(int)
	if !ok || userId == 0 {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	id := r.PathValue("id")
	userIDFromUrl, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

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

	// Call the service to get followers
	followers, err := h.service.GetFollowers(userIDFromUrl,userId, limit, offset, sort, orderBy, search)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}
func (h *FollowerHandler) GetFollowees(w http.ResponseWriter, r *http.Request) {	
	userId, ok := r.Context().Value(userIDKey).(interface{}).(int)
	if !ok || userId == 0 {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}
	id := r.PathValue("id")
	userIDFromUrl, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

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

	// Call the service to get followers
	followers, err := h.service.GetFollowees(userIDFromUrl,userId, limit, offset, sort, orderBy, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

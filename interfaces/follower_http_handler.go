package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bandvov/social-media-go/application"
)

type FollowerHandler struct {
	service *application.FollowerService
}

func NewFollowerHandler(service *application.FollowerService) *FollowerHandler {
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
		http.Error(w, "invalid follower ID", http.StatusBadRequest)
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
	// Parse the URL parameters to get the user ID
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Call the service to get followers
	followers, err := h.service.GetFollowers(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the followers into JSON and send as response
	response, err := json.Marshal(followers)
	if err != nil {
		http.Error(w, "Failed to marshal followers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

package interfaces

import (
	"activity-service/application"
	"encoding/json"
	"net/http"
	"strconv"
)

// Define keys for context
type contextKey string

const (
	userIDKey  contextKey = "userID"
	isAdminKey contextKey = "isAdmin"
)

// ActivityHandler handles HTTP requests for the activity service.
type ActivityHandler struct {
	service *application.ActivityService
}

// NewActivityHandler initializes a new handler.
func NewActivityHandler(service *application.ActivityService) *ActivityHandler {
	return &ActivityHandler{service: service}
}

// AddActivityEndpoint handles adding a new activity.
func (h *ActivityHandler) AddActivity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    string                 `json:"user_id"`
		Action    string                 `json:"action"`
		TargetID  string                 `json:"target_id"`
		EventData map[string]interface{} `json:"event_data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.AddActivity(req.UserID, req.Action, req.TargetID, req.EventData); err != nil {
		http.Error(w, "Failed to add activity", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetActivitiesEndpoint handles retrieving user activities.
func (h *ActivityHandler) GetActivities(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	activities, err := h.service.GetRecentActivities(userID, limit)
	if err != nil {
		http.Error(w, "Failed to fetch activities", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(activities)
}

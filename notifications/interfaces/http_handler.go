package interfaces

import (
	"encoding/json"
	"n/application"
	"net/http"
)

type NotificationHandler struct {
	service *application.NotificationService
}

func NewNotificationHandler(service *application.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// Send Notification Endpoint
func (h *NotificationHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID  string `json:"user_id"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.SendNotification(req.UserID, req.Message); err != nil {
		http.Error(w, "Failed to send notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Listen to Notifications (WebSocket / SSE)
func (h *NotificationHandler) ListenNotifications(w http.ResponseWriter, r *http.Request) {
	// Implementation for WebSocket / SSE
}

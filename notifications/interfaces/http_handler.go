package interfaces

import (
	"encoding/json"
	"fmt"
	"n/application"
	"n/domain"
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
	var req domain.Notification

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.SendNotification(req); err != nil {
		http.Error(w, "Failed to send notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Listen for notifications (SSE)
func (h *NotificationHandler) ListenNotifications(w http.ResponseWriter, r *http.Request) {
	// Get user_id from query parameters
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Subscribe to real-time notifications
	h.service.SubscribeToNotifications(userID, func(message string) {
		fmt.Fprintf(w, "data: %s\n\n", message)
		flusher.Flush()
	})

	// Keep connection open
	select {}
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Data []int `json:"data"`
	}

	// Parse and validate the request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "{\"message\": \"invalid request body\"}", http.StatusBadRequest)
		return
	}
	h.service.MarkAsRead(request.Data)
}

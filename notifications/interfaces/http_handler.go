package interfaces

import (
	"encoding/json"
	"fmt"
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

	// Fetch and send unsent messages
	messages, err := h.service.FetchUnsentMessages(userID)
	if err != nil {
		http.Error(w, "Failed to fetch unsent messages", http.StatusInternalServerError)
		return
	}
	messagesIds := make([]int, len(messages))
	for _, msg := range messages {
		fmt.Fprintf(w, "data: %s\n\n", msg.Message)
		flusher.Flush()
	}
	_ = h.service.MarkAsSent(messagesIds)

	// Subscribe to real-time notifications
	h.service.SubscribeToNotifications(userID, func(message string) {
		fmt.Fprintf(w, "data: %s\n\n", message)
		flusher.Flush()
	})

	// Keep connection open
	select {}
}

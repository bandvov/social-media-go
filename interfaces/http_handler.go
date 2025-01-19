package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/bandvov/social-media-go/application"
	"github.com/bandvov/social-media-go/domain"
)

type ReactionHandler struct {
	service application.ReactionServiceInterface
}

func NewReactionHandler(service application.ReactionServiceInterface) *ReactionHandler {
	return &ReactionHandler{service: service}
}

func (h *ReactionHandler) AddOrUpdateReaction(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value(userIDKey).(interface{}).(int)
	if !ok || userId == 0 {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}
	var reaction domain.Reaction
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.AddOrUpdateReaction(userId, reaction); err != nil {
		http.Error(w, "Failed to add or update reaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ReactionHandler) RemoveReaction(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	entityID := r.URL.Query().Get("entity_id")

	if userID == "" || entityID == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	if err := h.service.RemoveReaction(userID, entityID); err != nil {
		http.Error(w, "Failed to remove reaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

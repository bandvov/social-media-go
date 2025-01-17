package interfaces

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bandvov/social-media-go/application"
	"github.com/bandvov/social-media-go/domain"
)

type CommentHandler struct {
	service *application.CommentService
}

func NewCommentHandler(service *application.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Data *domain.Comment `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if !req.Data.IsValidAuthorId() || !req.Data.IsValidEntityId() || !req.Data.IsValidContent() {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}

	if err := h.service.AddComment(req.Data); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	entityID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	comments, err := h.service.GetComments(entityID)
	if err != nil {
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

package interfaces

import (
	"comments/application"
	"comments/domain"
	"comments/internal"
	"comments/utils"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type EntityIDsRequest struct {
	EntityIDs []int `json:"entity_ids"`
}

type CommentHandler struct {
	service *application.CommentService
	db      *sql.DB
	rdb     internal.RedisClient
}

func NewCommentHandler(service *application.CommentService, db *sql.DB, rdb internal.RedisClient) *CommentHandler {
	return &CommentHandler{service: service, db: db, rdb: rdb}
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

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.service.AddComment(ctx, req.Data); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CommentHandler) GetCommentsByEntityID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	entityID, err := strconv.Atoi(idStr)
	if err != nil || entityID <= 0 {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	targetUserId := query.Get("target_id")
	targetUserIDFromUrl, err := strconv.Atoi(targetUserId)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}
	p := utils.ParsePagination(r)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	comments, err := h.service.GetCommentsByEntityID(ctx, entityID, targetUserIDFromUrl, p)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":    comments,
		"hasMore": true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *CommentHandler) GetCommentsAndRepliesCount(w http.ResponseWriter, r *http.Request) {
	var request EntityIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	counts, err := h.service.GetCommentsAndRepliesCount(ctx, request.EntityIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counts)
}

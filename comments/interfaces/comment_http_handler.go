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
	"os"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
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
		Data domain.Comment `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message": "Invalid JSON request"}`, http.StatusBadRequest)
		return
	}
	if !req.Data.IsValidAuthorId() || !req.Data.IsValidEntityId() || !req.Data.IsValidContent() {
		http.Error(w, `{"message": "Missing required fields"}`, http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.service.AddComment(ctx, req.Data); err != nil {
		fmt.Println(err)
		http.Error(w, `{"message": "Failed to add comment"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CommentHandler) GetCommentsByEntityID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	entityID, err := strconv.Atoi(idStr)
	if err != nil || entityID <= 0 {
		http.Error(w, `{"message": "invalid entity ID"}`, http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	targetUserId := query.Get("target_id")
	targetUserIDFromUrl, err := strconv.Atoi(targetUserId)
	if err != nil {
		http.Error(w, `{"message": "invalid user ID"}`, http.StatusBadRequest)
		return
	}

	p := utils.ParsePagination(r)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var g errgroup.Group

	var comments []domain.Comment
	var counts []domain.CommentCount

	g.Go(func() error {
		var err error
		comments, err = h.service.GetCommentsByEntityID(ctx, entityID, targetUserIDFromUrl, p)
		if err != nil {
			fmt.Fprintln(os.Stdout, err)
			return fmt.Errorf("failed to get comments: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		counts, err = h.service.GetCommentsAndRepliesCount(ctx, []int{entityID})
		if err != nil {
			return fmt.Errorf("failed to get comments and replies count: %w", err)
		}
		return nil
	})

	// Wait for both operations to complete
	if err := g.Wait(); err != nil {
		http.Error(w, fmt.Sprintf(`{"message": %v}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Construct the response
	response := map[string]interface{}{
		"data":    comments,
		"hasMore": counts[0].CommentCount > len(comments),
	}

	// Set the content type and return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"message": "failed to encode response"}`, http.StatusInternalServerError)
	}
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

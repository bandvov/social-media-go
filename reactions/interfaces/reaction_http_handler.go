package interfaces

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reactions/application"
	"reactions/domain"
	"reactions/internal"
	"strconv"
	"time"
)

type ReactionHandler struct {
	service application.ReactionServiceInterface
	db      *sql.DB
	rdb     internal.RedisClient
}

func NewReactionHandler(service application.ReactionServiceInterface, db *sql.DB, rdb internal.RedisClient) *ReactionHandler {
	return &ReactionHandler{service: service, db: db, rdb: rdb}
}

func (h *ReactionHandler) AddOrUpdateReaction(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.Header.Get("X-User-Id"))
	if err != nil || userId == 0 {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}
	var reaction domain.Reaction
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.service.AddOrUpdateReaction(ctx, userId, reaction); err != nil {
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

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.service.RemoveReaction(ctx, userID, entityID); err != nil {
		http.Error(w, "Failed to remove reaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ReactionHandler) GetReactionsHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Data []domain.Entity `json:"data"`
	}
	// Decode request body to get entityIDs
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call the service method to get reactions
	reactions, err := s.service.GetReactions(r.Context(), req.Data)
	if err != nil {
		http.Error(w, "Failed to retrieve reactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reactions)
}

func (s *ReactionHandler) GetReactionsCount(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Data []domain.Entity `json:"data"`
	}

	// Decode the request body to get entityIDs
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call the service method to get the reactions count
	reactions, err := s.service.GetReactionsCount(r.Context(), req.Data)
	if err != nil {
		slog.Error("Failed to retrieve reactions count", "error", err)
		http.Error(w, "Failed to retrieve reactions count", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reactions); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (s *ReactionHandler) GetReactionStats(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Data []domain.Entity `json:"data"`
	}

	// Decode the request body to get entityIDs
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call the service method to get the reactions count
	reactions, err := s.service.GetReactionStats(r.Context(), req.Data)
	if err != nil {
		fmt.Fprintln(os.Stdout, "Failed to retrieve reactions stats", "error", err)
		http.Error(w, "Failed to retrieve reactions stats", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]domain.ReactionStat{"data": reactions}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *ReactionHandler) GetUserReactions(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Data []domain.Entity `json:"data"`
	}

	// Decode the request body to get entityIDs
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	userIDFromUrl, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	reactions, err := h.service.GetUserReactions(ctx, userIDFromUrl, req.Data)
	if err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
		http.Error(w, "Error fetching reaction", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(os.Stdout, "%+v", reactions)
	json.NewEncoder(w).Encode(map[string][]domain.Reaction{"data": reactions})
}

func (h *ReactionHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		slog.Warn("PostgreSQL health check failed", "error", err)
		http.Error(w, "Unhealthy", http.StatusServiceUnavailable)
		return
	}
	if err := h.rdb.Ping(ctx).Err(); err != nil {
		slog.Warn("Redis health check failed", "error", err)
		http.Error(w, "Unhealthy", http.StatusServiceUnavailable)
		return
	}

	slog.Info("Health check passed")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Ok")
}

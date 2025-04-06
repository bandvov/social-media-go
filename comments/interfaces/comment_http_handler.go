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
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.5.0"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
)

type EntityIDsRequest struct {
	EntityIDs []int `json:"entity_ids"`
}

type CommentHandler struct {
	service *application.CommentService
	db      *sql.DB
	rdb     internal.RedisClient
	tracer  trace.Tracer
}

func NewCommentHandler(service *application.CommentService, db *sql.DB, rdb internal.RedisClient, tracer trace.Tracer) *CommentHandler {
	return &CommentHandler{service: service, db: db, rdb: rdb, tracer: tracer}
}

func (h *CommentHandler) AddComment(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Parent span
	ctx, parentSpan := h.tracer.Start(ctx, "handler.AddComment")
	defer parentSpan.End()

	// Add an annotation when the request is received
	parentSpan.AddEvent("Request received", trace.WithAttributes(
		semconv.HTTPMethodKey.String(r.Method),
		semconv.HTTPURLKey.String(r.URL.String()),
		semconv.HTTPStatusCodeKey.Int(http.StatusOK),
	))

	var req struct {
		Data domain.Comment `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		parentSpan.SetStatus(codes.Error, err.Error())
		parentSpan.RecordError(err)
		http.Error(w, `{"message": "Invalid JSON request"}`, http.StatusBadRequest)
		return
	}

	if !req.Data.IsValidAuthorId() || !req.Data.IsValidEntityId() || !req.Data.IsValidContent() {
		parentSpan.SetStatus(codes.Error, "Missing required fields")
		http.Error(w, `{"message": "Missing required fields"}`, http.StatusBadRequest)
	}

	if err := h.service.AddComment(ctx, req.Data); err != nil {
		parentSpan.SetStatus(codes.Error, "Missing required fields")
		http.Error(w, `{"message": "Failed to add comment"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	parentSpan.SetStatus(codes.Ok, "success")

}

func (h *CommentHandler) GetCommentsByEntityID(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	// Parent span
	ctx, parentSpan := h.tracer.Start(ctx, "handler.GetCommentsByEntityID")
	defer parentSpan.End()
	// Add an annotation when the request is received
	parentSpan.AddEvent("Request received", trace.WithAttributes(
		semconv.HTTPMethodKey.String(r.Method),
		semconv.HTTPURLKey.String(r.URL.String()),
		semconv.HTTPStatusCodeKey.Int(http.StatusOK),
	))

	idStr := r.PathValue("id")
	entityID, err := strconv.Atoi(idStr)
	if err != nil || entityID <= 0 {
		parentSpan.AddEvent("validation-error", trace.WithAttributes(
			attribute.Bool("success", false),
			attribute.String("error", err.Error()),
		))
		http.Error(w, `{"message": "invalid entity ID"}`, http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	targetUserId := query.Get("target_id")
	targetUserIDFromUrl, err := strconv.Atoi(targetUserId)
	if err != nil {
		parentSpan.AddEvent("validation-error", trace.WithAttributes(
			attribute.Bool("success", false),
			attribute.String("error", err.Error()),
		))
		http.Error(w, `{"message": "invalid user ID"}`, http.StatusBadRequest)
		return
	}

	p := utils.ParsePagination(r)

	var g errgroup.Group

	var comments []domain.Comment
	var counts []domain.CommentCount

	g.Go(func() error {
		var err error
		comments, err = h.service.GetCommentsByEntityID(ctx, entityID, targetUserIDFromUrl, p)
		if err != nil {
			return fmt.Errorf("failed to get comments: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		counts, err = h.service.GetCommentsAndRepliesCount(ctx, []domain.Entity{
			{ID: entityID, Type: "post"},
		})
		if err != nil {
			return fmt.Errorf("failed to get comments and replies count: %w", err)
		}
		return nil
	})

	// Wait for both operations to complete
	if err := g.Wait(); err != nil {
		parentSpan.AddEvent("fetching-error", trace.WithAttributes(
			attribute.Bool("success", false),
			attribute.String("error", err.Error()),
		))
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
		parentSpan.AddEvent("encode-error", trace.WithAttributes(
			attribute.Bool("success", false),
			attribute.String("error", err.Error()),
		))
		http.Error(w, `{"message": "failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *CommentHandler) GetCommentsAndRepliesCount(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Data []domain.Entity `json:"data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"message": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	counts, err := h.service.GetCommentsAndRepliesCount(ctx, request.Data)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"message":%v}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counts)
}

func (h *CommentHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
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

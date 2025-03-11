package interfaces

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"posts/application"
	"posts/domain"
	"posts/internal"
	"posts/utils"
	"strconv"
	"time"
)

// Define keys for context
type contextKey string

const (
	userIDKey  contextKey = "userID"
	isAdminKey contextKey = "isAdmin"
)

type PostHTTPHandler struct {
	postService application.PostServiceInterface
	db          *sql.DB
	rdb         internal.RedisClient
}

func NewPostHTTPHandler(
	postService application.PostServiceInterface,
	rdb internal.RedisClient,

) *PostHTTPHandler {
	return &PostHTTPHandler{
		postService: postService,
		rdb:         rdb,
	}
}

func (p *PostHTTPHandler) CreatePost(w http.ResponseWriter, r *http.Request) {

	hv := r.Header.Get("user_id")

	authorID, err := strconv.Atoi(hv)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	var newPost struct {
		Data domain.CreatePostRequest `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if newPost.Data.Content == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	newPost.Data.AuthorID = authorID

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err = p.postService.CreatePost(ctx, &newPost.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
}

func (p *PostHTTPHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}
	if postID == 0 {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err = p.postService.DeletePost(ctx, postID)
	if err != nil {
		http.Error(w, "error deleting post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "post deleted successfully"})
}

func (p *PostHTTPHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}

	var post *domain.CreatePostRequest

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	err = p.postService.UpdatePost(ctx, postID, &domain.Post{
		Content: post.Content, Visibility: &post.Visibility, Tags: post.Tags, Pinned: post.Pinned,
	})

	if err != nil {
		http.Error(w, "error updating post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "post updated successfully"})
}

func (p *PostHTTPHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	targetUserIdStr := query.Get("target_id")
	targetUserId, err := strconv.Atoi(targetUserIdStr)
	if err != nil {
		http.Error(w, "invalid target id", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid post ID", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	post, err := p.postService.GetPostByID(ctx, postID, targetUserId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}
func (h *PostHTTPHandler) GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(os.Stdout, "GetPostsByUser")
	s := time.Now()
	query := r.URL.Query()
	targetUserIdStr := query.Get("target_id")
	targetUserId, err := strconv.Atoi(targetUserIdStr)
	if err != nil {
		http.Error(w, "invalid target id", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	authorIDFromUrl, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	p := utils.ParsePagination(r)
	// Move the logic to the service layer
	posts, postsCount, err := h.postService.GetPostsByUser(ctx, authorIDFromUrl, targetUserId, p)
	if err != nil {
		fmt.Fprintln(os.Stdout, err.Error())
		http.Error(w, "Failed to fetch posts", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"data":    posts,
		"hasMore": postsCount > p.Offset+p.Limit,
	}
	fmt.Println(time.Since(s))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PostHTTPHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
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
